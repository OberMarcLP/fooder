import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { LazyImage } from './LazyImage';

describe('LazyImage', () => {
  beforeEach(() => {
    // Reset IntersectionObserver mock
    global.IntersectionObserver = class IntersectionObserver {
      observe = vi.fn();
      disconnect = vi.fn();
      unobserve = vi.fn();
      takeRecords = vi.fn(() => []);

      constructor(callback: IntersectionObserverCallback) {
        // Automatically trigger callback to simulate image in view
        setTimeout(() => {
          callback(
            [
              {
                isIntersecting: true,
                target: document.createElement('img'),
              } as IntersectionObserverEntry,
            ],
            this as any
          );
        }, 0);
      }
    } as any;
  });

  it('renders an image', () => {
    render(<LazyImage src="test.jpg" alt="Test image" />);
    const img = screen.getByRole('img', { name: 'Test image' });
    expect(img).toBeInTheDocument();
  });

  it('shows loading skeleton initially', () => {
    const { container } = render(<LazyImage src="test.jpg" alt="Test image" />);
    const skeleton = container.querySelector('.animate-pulse');
    expect(skeleton).toBeInTheDocument();
  });

  it('applies className prop', () => {
    render(<LazyImage src="test.jpg" alt="Test image" className="custom-class" />);
    const img = screen.getByRole('img', { name: 'Test image' });
    expect(img).toHaveClass('custom-class');
  });

  it('calls onError when image fails to load', () => {
    const handleError = vi.fn();
    render(<LazyImage src="invalid.jpg" alt="Test image" onError={handleError} />);

    const img = screen.getByRole('img', { name: 'Test image' });
    const event = new Event('error');
    img.dispatchEvent(event);

    expect(handleError).toHaveBeenCalled();
  });

  it('loads image when in viewport', async () => {
    render(<LazyImage src="test.jpg" alt="Test image" />);

    await waitFor(() => {
      const img = screen.getByRole('img', { name: 'Test image' });
      expect(img).toHaveAttribute('src', 'test.jpg');
    });
  });

  it('uses IntersectionObserver', () => {
    const observeMock = vi.fn();
    global.IntersectionObserver = class IntersectionObserver {
      observe = observeMock;
      disconnect = vi.fn();
      unobserve = vi.fn();
      takeRecords = vi.fn(() => []);

      constructor() {}
    } as any;

    render(<LazyImage src="test.jpg" alt="Test image" />);

    // IntersectionObserver should be created and observe should be called
    expect(observeMock).toHaveBeenCalled();
  });
});
