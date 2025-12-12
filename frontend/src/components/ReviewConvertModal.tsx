import { useState } from 'react';
import { Modal } from './Modal';
import { StarRating } from './StarRating';

interface ReviewConvertModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: {
    foodRating: number;
    serviceRating: number;
    ambianceRating: number;
    comment: string;
    description: string;
  }) => void;
  restaurantName: string;
}

export function ReviewConvertModal({ isOpen, onClose, onSubmit, restaurantName }: ReviewConvertModalProps) {
  const [foodRating, setFoodRating] = useState(0);
  const [serviceRating, setServiceRating] = useState(0);
  const [ambianceRating, setAmbianceRating] = useState(0);
  const [comment, setComment] = useState('');
  const [description, setDescription] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (foodRating === 0 || serviceRating === 0 || ambianceRating === 0) {
      alert('Please provide all ratings');
      return;
    }

    onSubmit({
      foodRating,
      serviceRating,
      ambianceRating,
      comment,
      description,
    });

    // Reset form
    setFoodRating(0);
    setServiceRating(0);
    setAmbianceRating(0);
    setComment('');
    setDescription('');
  };

  const handleClose = () => {
    setFoodRating(0);
    setServiceRating(0);
    setAmbianceRating(0);
    setComment('');
    setDescription('');
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title={`Test & Review: ${restaurantName}`}>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Description (Optional)
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Add a description for this restaurant..."
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                     bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                     focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            rows={3}
          />
        </div>

        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            Your Review
          </h3>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Food Rating *
            </label>
            <StarRating rating={foodRating} onRatingChange={setFoodRating} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Service Rating *
            </label>
            <StarRating rating={serviceRating} onRatingChange={setServiceRating} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Ambiance Rating *
            </label>
            <StarRating rating={ambianceRating} onRatingChange={setAmbianceRating} />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Comment (Optional)
            </label>
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Share your experience..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                       bg-white dark:bg-gray-800 text-gray-900 dark:text-white
                       focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              rows={4}
            />
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
          <button
            type="button"
            onClick={handleClose}
            className="px-4 py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            Submit Review & Convert
          </button>
        </div>
      </form>
    </Modal>
  );
}
