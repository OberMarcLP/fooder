import { MapPin, Utensils, Tag } from 'lucide-react';
import { Restaurant } from '../services/api';
import { StarRating } from './StarRating';

interface RestaurantCardProps {
  restaurant: Restaurant;
  onClick: () => void;
}

export function RestaurantCard({ restaurant, onClick }: RestaurantCardProps) {
  return (
    <div
      onClick={onClick}
      className="card cursor-pointer hover:shadow-xl transition-shadow"
    >
      <h3 className="text-xl font-semibold mb-2">{restaurant.name}</h3>

      {restaurant.description && (
        <p className="text-gray-600 dark:text-gray-400 text-sm mb-3 line-clamp-2">
          {restaurant.description}
        </p>
      )}

      <div className="flex flex-wrap gap-2 mb-3">
        {restaurant.category && (
          <span className="inline-flex items-center gap-1 px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs rounded-full">
            <Tag className="w-3 h-3" />
            {restaurant.category.name}
          </span>
        )}
        {restaurant.food_types?.map((ft) => (
          <span key={ft.id} className="inline-flex items-center gap-1 px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 text-xs rounded-full">
            <Utensils className="w-3 h-3" />
            {ft.name}
          </span>
        ))}
      </div>

      {(restaurant.address || restaurant.distance !== undefined) && (
        <div className="flex items-start gap-2 text-sm text-gray-500 dark:text-gray-400 mb-3">
          <MapPin className="w-4 h-4 mt-0.5 flex-shrink-0" />
          <div className="flex-1">
            {restaurant.address && <span className="line-clamp-2">{restaurant.address}</span>}
            {restaurant.distance !== undefined && (
              <span className="text-blue-500 font-medium">
                {restaurant.address ? ' · ' : ''}{restaurant.distance.toFixed(1)} km away
              </span>
            )}
          </div>
        </div>
      )}

      {restaurant.avg_rating && (
        <div className="flex items-center gap-2 pt-3 border-t border-gray-200 dark:border-gray-700">
          <StarRating rating={Math.round(restaurant.avg_rating.overall)} readonly size="sm" />
          <span className="text-sm text-gray-600 dark:text-gray-400">
            {restaurant.avg_rating.overall.toFixed(1)} ({restaurant.avg_rating.count} reviews)
          </span>
        </div>
      )}
    </div>
  );
}
