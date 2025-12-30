import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { StarRating } from './StarRating';

describe('StarRating', () => {
  it('renders the correct number of stars', () => {
    const { container } = render(<StarRating rating={3} />);
    const stars = container.querySelectorAll('svg');
    expect(stars).toHaveLength(5);
  });

  it('displays the correct rating value', () => {
    render(<StarRating rating={4} />);
    expect(screen.getByText('4.0')).toBeInTheDocument();
  });

  it('renders filled stars based on rating', () => {
    const { container } = render(<StarRating rating={3} />);
    const filledStars = container.querySelectorAll('.text-yellow-400');
    expect(filledStars.length).toBeGreaterThan(0);
  });

  it('calls onChange when a star is clicked in interactive mode', () => {
    const handleChange = vi.fn();
    const { container } = render(<StarRating rating={0} onChange={handleChange} />);

    const stars = container.querySelectorAll('svg');
    fireEvent.click(stars[2]); // Click third star

    expect(handleChange).toHaveBeenCalledWith(3);
  });

  it('does not call onChange when readonly', () => {
    const handleChange = vi.fn();
    const { container } = render(<StarRating rating={3} onChange={handleChange} readonly />);

    const stars = container.querySelectorAll('svg');
    fireEvent.click(stars[4]);

    expect(handleChange).not.toHaveBeenCalled();
  });

  it('shows hover effect in interactive mode', () => {
    const handleChange = vi.fn();
    const { container } = render(<StarRating rating={0} onChange={handleChange} />);

    const stars = container.querySelectorAll('svg');
    fireEvent.mouseEnter(stars[3]);

    // Verify hover state is applied
    expect(container.querySelector('.cursor-pointer')).toBeInTheDocument();
  });

  it('handles half ratings correctly', () => {
    render(<StarRating rating={3.5} />);
    expect(screen.getByText('3.5')).toBeInTheDocument();
  });

  it('handles zero rating', () => {
    render(<StarRating rating={0} />);
    expect(screen.getByText('0.0')).toBeInTheDocument();
  });

  it('handles maximum rating', () => {
    render(<StarRating rating={5} />);
    expect(screen.getByText('5.0')).toBeInTheDocument();
  });
});
