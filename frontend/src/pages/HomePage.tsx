import { useState, useEffect, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Plus, Loader2 } from 'lucide-react';
import { Restaurant, CreateRestaurantData, CreateSuggestionData, RestaurantFilters, getRestaurants, createSuggestion, updateRestaurant, deleteRestaurant, getRestaurant, convertSuggestion, deleteSuggestion } from '../services/api';
import { RestaurantCard } from '../components/RestaurantCard';
import { RestaurantForm } from '../components/RestaurantForm';
import { SuggestionForm } from '../components/SuggestionForm';
import { ReviewConvertModal } from '../components/ReviewConvertModal';
import { Modal } from '../components/Modal';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { AlertDialog } from '../components/AlertDialog';
import { RestaurantDetail } from './RestaurantDetail';

interface HomePageProps {
  filters: RestaurantFilters;
}

export function HomePage({ filters }: HomePageProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddSuggestionModal, setShowAddSuggestionModal] = useState(false);
  const [selectedRestaurant, setSelectedRestaurant] = useState<Restaurant | null>(null);
  const [editingRestaurant, setEditingRestaurant] = useState<Restaurant | null>(null);
  const [reviewingSuggestion, setReviewingSuggestion] = useState<Restaurant | null>(null);
  const [rejectingRestaurant, setRejectingRestaurant] = useState<Restaurant | null>(null);
  const [deletingRestaurant, setDeletingRestaurant] = useState<Restaurant | null>(null);
  const [alertMessage, setAlertMessage] = useState<string>('');

  const fetchRestaurants = useCallback(async () => {
    try {
      const data = await getRestaurants(filters);
      setRestaurants(data);
    } catch (error) {
      console.error('Failed to fetch restaurants:', error);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchRestaurants();
  }, [fetchRestaurants]);

  // Handle restaurant selection from global search
  useEffect(() => {
    const state = location.state as { restaurantId?: number } | null;
    if (state?.restaurantId) {
      const restaurantId = state.restaurantId;
      const fetchSelectedRestaurant = async () => {
        try {
          const restaurant = await getRestaurant(restaurantId);
          setSelectedRestaurant(restaurant);
          // Clear the location state
          navigate(location.pathname, { replace: true });
        } catch (error) {
          console.error('Failed to fetch restaurant:', error);
        }
      };
      fetchSelectedRestaurant();
    }
  }, [location, navigate]);

  const handleCreateSuggestion = async (data: CreateSuggestionData) => {
    try {
      await createSuggestion(data);
      setShowAddSuggestionModal(false);
      fetchRestaurants();
    } catch (error: any) {
      console.error('Failed to create suggestion:', error);
      if (error.message && (error.message.includes('already exists') || error.message.includes('duplicate'))) {
        setAlertMessage('This restaurant already exists in the database. Please search for it instead.');
      } else {
        setAlertMessage('Failed to create suggestion. Please try again.');
      }
    }
  };

  const handleUpdate = async (data: CreateRestaurantData) => {
    if (!editingRestaurant) return;
    try {
      await updateRestaurant(editingRestaurant.id, data);
      setEditingRestaurant(null);
      fetchRestaurants();
    } catch (error) {
      console.error('Failed to update restaurant:', error);
    }
  };

  const handleDelete = async () => {
    setDeletingRestaurant(selectedRestaurant);
  };

  const confirmDelete = async () => {
    if (!deletingRestaurant) return;
    try {
      await deleteRestaurant(deletingRestaurant.id);
      setSelectedRestaurant(null);
      setDeletingRestaurant(null);
      fetchRestaurants();
    } catch (error) {
      console.error('Failed to delete restaurant:', error);
    }
  };

  const handleReviewSuggestion = (restaurant: Restaurant) => {
    setReviewingSuggestion(restaurant);
  };

  const handleRejectSuggestion = async (restaurant: Restaurant) => {
    setRejectingRestaurant(restaurant);
  };

  const confirmReject = async () => {
    if (!rejectingRestaurant) return;
    try {
      if (rejectingRestaurant.suggestion_id) {
        await deleteSuggestion(rejectingRestaurant.suggestion_id);
        setRejectingRestaurant(null);
        fetchRestaurants();
      }
    } catch (error) {
      console.error('Failed to reject suggestion:', error);
    }
  };

  const handleConvertSuggestion = async (data: {
    foodRating: number;
    serviceRating: number;
    ambianceRating: number;
    comment: string;
    description: string;
  }) => {
    if (!reviewingSuggestion?.suggestion_id) return;
    try {
      await convertSuggestion(reviewingSuggestion.suggestion_id, {
        description: data.description || undefined,
        category_id: reviewingSuggestion.category?.id || undefined,
        food_rating: data.foodRating,
        service_rating: data.serviceRating,
        ambiance_rating: data.ambianceRating,
        comment: data.comment || undefined,
      });
      setReviewingSuggestion(null);
      fetchRestaurants();
    } catch (error: any) {
      console.error('Failed to convert suggestion:', error);
      if (error.message && (error.message.includes('already exists') || error.message.includes('duplicate'))) {
        setAlertMessage('This restaurant already exists in the database. The suggestion will be rejected.');
        if (reviewingSuggestion?.suggestion_id) {
          await deleteSuggestion(reviewingSuggestion.suggestion_id);
          setReviewingSuggestion(null);
          fetchRestaurants();
        }
      } else {
        setAlertMessage('Failed to convert suggestion. Please try again.');
      }
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Restaurants</h1>
        <button onClick={() => setShowAddSuggestionModal(true)} className="btn btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add Suggestion
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
        </div>
      ) : restaurants.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-gray-400 mb-4">
            {Object.keys(filters).length > 0 ? 'No restaurants match your filters.' : 'No restaurants yet.'}
          </p>
          {Object.keys(filters).length === 0 && (
            <button onClick={() => setShowAddSuggestionModal(true)} className="btn btn-primary">
              Add your first suggestion
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {restaurants.map((restaurant) => (
            <RestaurantCard
              key={restaurant.id}
              restaurant={restaurant}
              onClick={() => setSelectedRestaurant(restaurant)}
              onReview={handleReviewSuggestion}
              onReject={handleRejectSuggestion}
            />
          ))}
        </div>
      )}

      <Modal isOpen={showAddSuggestionModal} onClose={() => setShowAddSuggestionModal(false)} title="Add Suggestion">
        <SuggestionForm onSubmit={handleCreateSuggestion} onCancel={() => setShowAddSuggestionModal(false)} />
      </Modal>

      <Modal
        isOpen={editingRestaurant !== null}
        onClose={() => setEditingRestaurant(null)}
        title="Edit Restaurant"
      >
        {editingRestaurant && (
          <RestaurantForm
            restaurant={editingRestaurant}
            onSubmit={handleUpdate}
            onCancel={() => setEditingRestaurant(null)}
          />
        )}
      </Modal>

      <Modal
        isOpen={selectedRestaurant !== null}
        onClose={() => setSelectedRestaurant(null)}
        title={selectedRestaurant?.name || ''}
      >
        {selectedRestaurant && (
          <RestaurantDetail
            restaurant={selectedRestaurant}
            onEdit={() => {
              setEditingRestaurant(selectedRestaurant);
              setSelectedRestaurant(null);
            }}
            onDelete={handleDelete}
            onRatingAdded={fetchRestaurants}
          />
        )}
      </Modal>

      {reviewingSuggestion && (
        <ReviewConvertModal
          isOpen={true}
          onClose={() => setReviewingSuggestion(null)}
          onSubmit={handleConvertSuggestion}
          restaurantName={reviewingSuggestion.name}
        />
      )}

      <ConfirmDialog
        isOpen={rejectingRestaurant !== null}
        onClose={() => setRejectingRestaurant(null)}
        onConfirm={confirmReject}
        title="Reject Suggestion"
        message={`Are you sure you want to reject "${rejectingRestaurant?.name}"?`}
        confirmText="Reject"
        cancelText="Cancel"
        confirmClassName="bg-red-600 hover:bg-red-700 text-white"
      />

      <ConfirmDialog
        isOpen={deletingRestaurant !== null}
        onClose={() => setDeletingRestaurant(null)}
        onConfirm={confirmDelete}
        title="Delete Restaurant"
        message={`Are you sure you want to delete "${deletingRestaurant?.name}"? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        confirmClassName="bg-red-600 hover:bg-red-700 text-white"
      />

      <AlertDialog
        isOpen={alertMessage !== ''}
        onClose={() => setAlertMessage('')}
        message={alertMessage}
      />
    </div>
  );
}
