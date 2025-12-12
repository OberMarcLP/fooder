import { useState, useEffect } from 'react';
import { MapPin, X, Loader2, Navigation } from 'lucide-react';
import { Category, FoodType, RestaurantFilters, geocodeCities, GooglePlaceResult } from '../services/api';

interface SearchFiltersProps {
  categories: Category[];
  foodTypes: FoodType[];
  filters: RestaurantFilters;
  onFiltersChange: (filters: RestaurantFilters) => void;
}

export function SearchFilters({ categories, foodTypes, filters, onFiltersChange }: SearchFiltersProps) {
  const [locationSearch, setLocationSearch] = useState('');
  const [locationResults, setLocationResults] = useState<GooglePlaceResult[]>([]);
  const [searchingLocation, setSearchingLocation] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState<{ name: string; lat: number; lng: number } | null>(null);
  const [gettingCurrentLocation, setGettingCurrentLocation] = useState(false);

  const hasActiveFilters = filters.category_id || (filters.food_type_ids && filters.food_type_ids.length > 0) || filters.radius;

  // Search for locations
  useEffect(() => {
    if (locationSearch.length < 2) {
      setLocationResults([]);
      return;
    }

    const timer = setTimeout(async () => {
      setSearchingLocation(true);
      try {
        const results = await geocodeCities(locationSearch);
        setLocationResults(results);
      } catch (error) {
        console.error('Failed to search locations:', error);
      } finally {
        setSearchingLocation(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [locationSearch]);

  const handleCategoryChange = (categoryId: number | undefined) => {
    onFiltersChange({ ...filters, category_id: categoryId });
  };

  const handleFoodTypeToggle = (foodTypeId: number) => {
    const currentIds = filters.food_type_ids || [];
    const newIds = currentIds.includes(foodTypeId)
      ? currentIds.filter(id => id !== foodTypeId)
      : [...currentIds, foodTypeId];
    onFiltersChange({ ...filters, food_type_ids: newIds.length > 0 ? newIds : undefined });
  };

  const handleRadiusChange = (radius: number | undefined) => {
    if (radius && selectedLocation) {
      onFiltersChange({
        ...filters,
        radius,
        lat: selectedLocation.lat,
        lng: selectedLocation.lng,
      });
    } else {
      onFiltersChange({
        ...filters,
        radius: undefined,
        lat: undefined,
        lng: undefined,
      });
    }
  };

  const handleLocationSelect = (place: GooglePlaceResult) => {
    setSelectedLocation({ name: place.name, lat: place.latitude, lng: place.longitude });
    setLocationSearch('');
    setLocationResults([]);
    if (filters.radius) {
      onFiltersChange({
        ...filters,
        lat: place.latitude,
        lng: place.longitude,
      });
    }
  };

  const handleUseCurrentLocation = () => {
    if (!navigator.geolocation) {
      alert('Geolocation is not supported by your browser');
      return;
    }

    setGettingCurrentLocation(true);
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        setSelectedLocation({ name: 'Current Location', lat: latitude, lng: longitude });
        if (filters.radius) {
          onFiltersChange({
            ...filters,
            lat: latitude,
            lng: longitude,
          });
        }
        setGettingCurrentLocation(false);
      },
      (error) => {
        console.error('Error getting location:', error);
        alert('Unable to get your location. Please search for a location instead.');
        setGettingCurrentLocation(false);
      }
    );
  };

  const clearLocationFilter = () => {
    setSelectedLocation(null);
    onFiltersChange({
      ...filters,
      lat: undefined,
      lng: undefined,
      radius: undefined,
    });
  };

  const clearAllFilters = () => {
    setSelectedLocation(null);
    onFiltersChange({});
  };

  return (
    <div className="card mb-6">
      {hasActiveFilters && (
        <div className="flex items-center justify-end mb-4">
          <button
            onClick={clearAllFilters}
            className="text-sm text-gray-500 hover:text-red-500 transition-colors"
          >
            Clear all
          </button>
        </div>
      )}

      <div className="space-y-4">
          {/* Category Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Category
            </label>
            <select
              value={filters.category_id || ''}
              onChange={(e) => handleCategoryChange(e.target.value ? Number(e.target.value) : undefined)}
              className="input"
            >
              <option value="">All Categories</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>{cat.name}</option>
              ))}
            </select>
          </div>

          {/* Food Types Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Food Types
            </label>
            <div className="flex flex-wrap gap-2">
              {foodTypes.map((ft) => (
                <button
                  key={ft.id}
                  onClick={() => handleFoodTypeToggle(ft.id)}
                  className={`px-3 py-1 rounded-full text-sm transition-colors ${
                    filters.food_type_ids?.includes(ft.id)
                      ? 'bg-green-500 text-white'
                      : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
                  }`}
                >
                  {ft.name}
                </button>
              ))}
            </div>
          </div>

          {/* Location Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Location & Radius
            </label>

            {selectedLocation ? (
              <div className="flex items-center gap-2 mb-3 p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
                <MapPin className="w-4 h-4 text-blue-500" />
                <span className="text-sm text-gray-700 dark:text-gray-300 flex-1">
                  {selectedLocation.name}
                </span>
                <button
                  onClick={clearLocationFilter}
                  className="text-gray-400 hover:text-red-500"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            ) : (
              <div className="space-y-2 mb-3">
                <div className="relative">
                  <input
                    type="text"
                    value={locationSearch}
                    onChange={(e) => setLocationSearch(e.target.value)}
                    placeholder="Search city..."
                    className="input !pl-10"
                  />
                  <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
                  {searchingLocation && (
                    <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 animate-spin" />
                  )}
                </div>

                {locationResults.length > 0 && (
                  <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg max-h-48 overflow-y-auto">
                    {locationResults.map((place) => (
                      <button
                        key={place.place_id}
                        onClick={() => handleLocationSelect(place)}
                        className="w-full px-3 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700 text-sm"
                      >
                        <div className="font-medium text-gray-900 dark:text-white">{place.name}</div>
                        <div className="text-gray-500 dark:text-gray-400 text-xs">{place.address}</div>
                      </button>
                    ))}
                  </div>
                )}

                <button
                  onClick={handleUseCurrentLocation}
                  disabled={gettingCurrentLocation}
                  className="flex items-center gap-2 text-sm text-blue-500 hover:text-blue-600 disabled:opacity-50"
                >
                  {gettingCurrentLocation ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <Navigation className="w-4 h-4" />
                  )}
                  Use my current location
                </button>
              </div>
            )}

            {/* Radius Selection */}
            <div className="flex flex-wrap gap-2">
              {[1, 5, 10, 25, 50].map((km) => (
                <button
                  key={km}
                  onClick={() => handleRadiusChange(filters.radius === km ? undefined : km)}
                  disabled={!selectedLocation}
                  className={`px-3 py-1 rounded-full text-sm transition-colors ${
                    filters.radius === km
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed'
                  }`}
                >
                  {km} km
                </button>
              ))}
            </div>
            {!selectedLocation && (
              <p className="text-xs text-gray-500 mt-1">Select a location first to filter by radius</p>
            )}
          </div>
      </div>
    </div>
  );
}
