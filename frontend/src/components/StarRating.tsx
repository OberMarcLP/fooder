import { Star } from 'lucide-react';

interface StarRatingProps {
  rating: number;
  onRatingChange?: (rating: number) => void;
  readonly?: boolean;
  size?: 'sm' | 'md' | 'lg';
}

export function StarRating({
  rating,
  onRatingChange,
  readonly = false,
  size = 'md',
}: StarRatingProps) {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
  };

  return (
    <div className="flex gap-1">
      {[1, 2, 3, 4, 5].map((star) => (
        <button
          key={star}
          type="button"
          onClick={() => !readonly && onRatingChange?.(star)}
          disabled={readonly}
          className={`${readonly ? 'cursor-default' : 'cursor-pointer hover:scale-125'} transition-all duration-200`}
        >
          <Star
            className={`${sizeClasses[size]} transition-all duration-200 ${
              star <= rating
                ? 'fill-yellow-400 text-yellow-400 drop-shadow-[0_0_8px_rgba(251,191,36,0.6)] hover:drop-shadow-[0_0_12px_rgba(251,191,36,0.8)]'
                : 'fill-none text-gray-300 dark:text-gray-600'
            }`}
          />
        </button>
      ))}
    </div>
  );
}
