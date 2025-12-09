import { useState, useEffect } from 'react';
import { Plus, Loader2 } from 'lucide-react';
import { Restaurant, CreateRestaurantData, getRestaurants, createRestaurant, updateRestaurant, deleteRestaurant } from '../services/api';
import { RestaurantCard } from '../components/RestaurantCard';
import { RestaurantForm } from '../components/RestaurantForm';
import { Modal } from '../components/Modal';
import { RestaurantDetail } from './RestaurantDetail';

export function HomePage() {
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [selectedRestaurant, setSelectedRestaurant] = useState<Restaurant | null>(null);
  const [editingRestaurant, setEditingRestaurant] = useState<Restaurant | null>(null);

  const fetchRestaurants = async () => {
    try {
      const data = await getRestaurants();
      setRestaurants(data);
    } catch (error) {
      console.error('Failed to fetch restaurants:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRestaurants();
  }, []);

  const handleCreate = async (data: CreateRestaurantData) => {
    try {
      await createRestaurant(data);
      setShowAddModal(false);
      fetchRestaurants();
    } catch (error) {
      console.error('Failed to create restaurant:', error);
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
        <button onClick={() => setShowAddModal(true)} className="btn btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add Restaurant
        </button>
      </div>

      {restaurants.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-gray-400 mb-4">No restaurants yet.</p>
          <button onClick={() => setShowAddModal(true)} className="btn btn-primary">
            Add your first restaurant
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {restaurants.map((restaurant) => (
            <RestaurantCard
              key={restaurant.id}
              restaurant={restaurant}
              onClick={() => setSelectedRestaurant(restaurant)}
            />
          ))}
        </div>
      )}

      <Modal isOpen={showAddModal} onClose={() => setShowAddModal(false)} title="Add Restaurant">
        <RestaurantForm onSubmit={handleCreate} onCancel={() => setShowAddModal(false)} />
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
    </div>
  );
}
