const API_URL = import.meta.env.VITE_API_URL || '';

export interface Category {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface FoodType {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface AvgRating {
  food: number;
  service: number;
  ambiance: number;
  overall: number;
  count: number;
}

export interface Restaurant {
  id: number;
  name: string;
  description: string | null;
  address: string | null;
  phone: string | null;
  website: string | null;
  latitude: number | null;
  longitude: number | null;
  google_place_id: string | null;
  category_id: number | null;
  category?: Category;
  food_types?: FoodType[];
  avg_rating?: AvgRating;
  created_at: string;
  updated_at: string;
}

export interface CreateRestaurantData {
  name: string;
  description?: string | null;
  address?: string | null;
  phone?: string | null;
  website?: string | null;
  latitude?: number | null;
  longitude?: number | null;
  google_place_id?: string | null;
  category_id?: number | null;
  food_type_ids?: number[];
}

export interface Rating {
  id: number;
  restaurant_id: number;
  food_rating: number;
  service_rating: number;
  ambiance_rating: number;
  comment: string | null;
  created_at: string;
}

export interface GooglePlaceResult {
  place_id: string;
  name: string;
  address: string;
  phone?: string;
  website?: string;
  latitude: number;
  longitude: number;
}

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}/api${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'An error occurred');
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

// Categories
export const getCategories = () => fetchApi<Category[]>('/categories');
export const createCategory = (name: string) =>
  fetchApi<Category>('/categories', {
    method: 'POST',
    body: JSON.stringify({ name }),
  });
export const updateCategory = (id: number, name: string) =>
  fetchApi<Category>(`/categories/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ name }),
  });
export const deleteCategory = (id: number) =>
  fetchApi<void>(`/categories/${id}`, { method: 'DELETE' });

// Food Types
export const getFoodTypes = () => fetchApi<FoodType[]>('/food-types');
export const createFoodType = (name: string) =>
  fetchApi<FoodType>('/food-types', {
    method: 'POST',
    body: JSON.stringify({ name }),
  });
export const updateFoodType = (id: number, name: string) =>
  fetchApi<FoodType>(`/food-types/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ name }),
  });
export const deleteFoodType = (id: number) =>
  fetchApi<void>(`/food-types/${id}`, { method: 'DELETE' });

// Restaurants
export const getRestaurants = () => fetchApi<Restaurant[]>('/restaurants');
export const getRestaurant = (id: number) => fetchApi<Restaurant>(`/restaurants/${id}`);
export const createRestaurant = (data: CreateRestaurantData) =>
  fetchApi<Restaurant>('/restaurants', {
    method: 'POST',
    body: JSON.stringify(data),
  });
export const updateRestaurant = (id: number, data: CreateRestaurantData) =>
  fetchApi<Restaurant>(`/restaurants/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
export const deleteRestaurant = (id: number) =>
  fetchApi<void>(`/restaurants/${id}`, { method: 'DELETE' });

// Ratings
export const getRatings = (restaurantId: number) =>
  fetchApi<Rating[]>(`/restaurants/${restaurantId}/ratings`);
export const createRating = (data: {
  restaurant_id: number;
  food_rating: number;
  service_rating: number;
  ambiance_rating: number;
  comment?: string;
}) =>
  fetchApi<Rating>('/ratings', {
    method: 'POST',
    body: JSON.stringify(data),
  });
export const deleteRating = (id: number) =>
  fetchApi<void>(`/ratings/${id}`, { method: 'DELETE' });

// Google Maps
export const searchPlaces = (query: string) =>
  fetchApi<GooglePlaceResult[]>(`/places/search?q=${encodeURIComponent(query)}`);
export const getPlaceDetails = (placeId: string) =>
  fetchApi<GooglePlaceResult>(`/places/${placeId}`);
