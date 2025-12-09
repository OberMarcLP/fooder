import { useState } from 'react';
import { StarRating } from './StarRating';

interface RatingFormProps {
  onSubmit: (data: {
    food_rating: number;
    service_rating: number;
    ambiance_rating: number;
    comment?: string;
  }) => void;
  onCancel: () => void;
}

export function RatingForm({ onSubmit, onCancel }: RatingFormProps) {
  const [foodRating, setFoodRating] = useState(0);
  const [serviceRating, setServiceRating] = useState(0);
  const [ambianceRating, setAmbianceRating] = useState(0);
  const [comment, setComment] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (foodRating === 0 || serviceRating === 0 || ambianceRating === 0) {
      return;
    }
    onSubmit({
      food_rating: foodRating,
      service_rating: serviceRating,
      ambiance_rating: ambianceRating,
      comment: comment || undefined,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label className="label mb-2">Food Rating *</label>
        <StarRating rating={foodRating} onRatingChange={setFoodRating} size="lg" />
      </div>

      <div>
        <label className="label mb-2">Service Rating *</label>
        <StarRating rating={serviceRating} onRatingChange={setServiceRating} size="lg" />
      </div>

      <div>
        <label className="label mb-2">Ambiance Rating *</label>
        <StarRating rating={ambianceRating} onRatingChange={setAmbianceRating} size="lg" />
      </div>

      <div>
        <label className="label">Comment (optional)</label>
        <textarea
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          className="input min-h-[100px]"
          placeholder="Share your experience..."
          rows={3}
        />
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={foodRating === 0 || serviceRating === 0 || ambianceRating === 0}
          className="btn btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Submit Rating
        </button>
        <button type="button" onClick={onCancel} className="btn btn-secondary">
          Cancel
        </button>
      </div>
    </form>
  );
}
