import { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Loader2, Utensils } from 'lucide-react';
import { FoodType, getFoodTypes, createFoodType, updateFoodType, deleteFoodType } from '../services/api';

export function FoodTypesPage() {
  const [foodTypes, setFoodTypes] = useState<FoodType[]>([]);
  const [loading, setLoading] = useState(true);
  const [newName, setNewName] = useState('');
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState('');

  const fetchFoodTypes = async () => {
    try {
      const data = await getFoodTypes();
      setFoodTypes(data);
    } catch (error) {
      console.error('Failed to fetch food types:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFoodTypes();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim()) return;
    try {
      await createFoodType(newName);
      setNewName('');
      fetchFoodTypes();
    } catch (error) {
      console.error('Failed to create food type:', error);
    }
  };

  const handleUpdate = async (id: number) => {
    if (!editName.trim()) return;
    try {
      await updateFoodType(id, editName);
      setEditingId(null);
      setEditName('');
      fetchFoodTypes();
    } catch (error) {
      console.error('Failed to update food type:', error);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this food type?')) return;
    try {
      await deleteFoodType(id);
      fetchFoodTypes();
    } catch (error) {
      console.error('Failed to delete food type:', error);
    }
  };

  const startEdit = (foodType: FoodType) => {
    setEditingId(foodType.id);
    setEditName(foodType.name);
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
      <div className="flex items-center gap-3 mb-6">
        <Utensils className="w-6 h-6 text-green-500" />
        <h1 className="text-2xl font-bold">Food Types</h1>
      </div>

      <form onSubmit={handleCreate} className="flex gap-2 mb-6">
        <input
          type="text"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          placeholder="New food type name..."
          className="input flex-1"
        />
        <button type="submit" className="btn btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add
        </button>
      </form>

      {foodTypes.length === 0 ? (
        <p className="text-center text-gray-500 dark:text-gray-400 py-8">
          No food types yet. Add your first food type above.
        </p>
      ) : (
        <div className="space-y-2">
          {foodTypes.map((foodType) => (
            <div key={foodType.id} className="card flex items-center justify-between">
              {editingId === foodType.id ? (
                <div className="flex items-center gap-2 flex-1">
                  <input
                    type="text"
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    className="input flex-1"
                    autoFocus
                  />
                  <button
                    onClick={() => handleUpdate(foodType.id)}
                    className="btn btn-primary text-sm"
                  >
                    Save
                  </button>
                  <button
                    onClick={() => setEditingId(null)}
                    className="btn btn-secondary text-sm"
                  >
                    Cancel
                  </button>
                </div>
              ) : (
                <>
                  <span className="font-medium">{foodType.name}</span>
                  <div className="flex gap-2">
                    <button
                      onClick={() => startEdit(foodType)}
                      className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg"
                    >
                      <Edit className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(foodType.id)}
                      className="p-2 hover:bg-red-100 dark:hover:bg-red-900/30 text-red-600 dark:text-red-400 rounded-lg"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
