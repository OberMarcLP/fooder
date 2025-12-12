import { MapPin, Utensils, Tag, Lightbulb, CheckCircle, XCircle } from 'lucide-react';
import { Restaurant } from '../services/api';
import { StarRating } from './StarRating';

interface RestaurantCardProps {
  restaurant: Restaurant;
  onClick: () => void;
  onReview?: (restaurant: Restaurant) => void;
  onReject?: (restaurant: Restaurant) => void;
}

export function RestaurantCard({ restaurant, onClick, onReview, onReject }: RestaurantCardProps) {
  return (
    <div
      className="card hover:shadow-xl transition-shadow"
    >
      <div onClick={onClick} className="cursor-pointer">
        <div className="flex items-start justify-between gap-2 mb-2">
          <h3 className="text-xl font-semibold flex-1">{restaurant.name}</h3>
          {restaurant.is_suggestion && (
            <span className="inline-flex items-center gap-1 px-2 py-1 bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 text-xs rounded-full flex-shrink-0">
              <Lightbulb className="w-3 h-3" />
              Suggestion
            </span>
          )}
        </div>

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

        {restaurant.is_suggestion ? (
          <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
            <p className="text-sm text-yellow-700 dark:text-yellow-400 italic">
              Not yet rated - Try it and add your review!
            </p>
          </div>
        ) : restaurant.avg_rating && (
          <div className="flex items-center gap-2 pt-3 border-t border-gray-200 dark:border-gray-700">
            <StarRating rating={Math.round(restaurant.avg_rating.overall)} readonly size="sm" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {restaurant.avg_rating.overall.toFixed(1)} ({restaurant.avg_rating.count} reviews)
            </span>
          </div>
        )}
      </div>

      {restaurant.is_suggestion && onReview && onReject && (
        <div className="flex gap-2 mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={(e) => {
              e.stopPropagation();
              onReview(restaurant);
            }}
            className="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
          >
            <CheckCircle className="w-4 h-4" />
            Review
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              onReject(restaurant);
            }}
            className="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            <XCircle className="w-4 h-4" />
            Reject
          </button>
        </div>
      )}
    </div>
  );
}
