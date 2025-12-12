import { useState, useEffect, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Plus, Loader2 } from 'lucide-react';
import { Restaurant, CreateRestaurantData, CreateSuggestionData, RestaurantFilters, getRestaurants, createSuggestion, updateRestaurant, deleteRestaurant, getRestaurant, convertSuggestion, deleteSuggestion } from '../services/api';
import { RestaurantCard } from '../components/RestaurantCard';
import { RestaurantForm } from '../components/RestaurantForm';
import { SuggestionForm } from '../components/SuggestionForm';
import { ReviewConvertModal } from '../components/ReviewConvertModal';
import { Modal } from '../components/Modal';
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
    } catch (error) {
      console.error('Failed to create suggestion:', error);
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

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this restaurant?')) return;
    try {
      await deleteRestaurant(id);
      setSelectedRestaurant(null);
      fetchRestaurants();
    } catch (error) {
      console.error('Failed to delete restaurant:', error);
    }
  };

  const handleReviewSuggestion = (restaurant: Restaurant) => {
    setReviewingSuggestion(restaurant);
  };

  const handleRejectSuggestion = async (restaurant: Restaurant) => {
    if (!confirm(`Are you sure you want to reject "${restaurant.name}"?`)) return;
    try {
      if (restaurant.suggestion_id) {
        await deleteSuggestion(restaurant.suggestion_id);
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
    } catch (error) {
      console.error('Failed to convert suggestion:', error);
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
            onDelete={() => handleDelete(selectedRestaurant.id)}
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
    </div>
  );
}
