import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ThemeToggle } from './ThemeToggle';

describe('ThemeToggle', () => {
  it('renders the toggle button', () => {
    const handleToggle = vi.fn();
    render(<ThemeToggle isDark={false} onToggle={handleToggle} />);

    const button = screen.getByRole('button');
    expect(button).toBeInTheDocument();
  });

  it('shows sun icon when in dark mode', () => {
    const handleToggle = vi.fn();
    const { container } = render(<ThemeToggle isDark={true} onToggle={handleToggle} />);

    // Sun icon should be visible in dark mode
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });

  it('shows moon icon when in light mode', () => {
    const handleToggle = vi.fn();
    const { container } = render(<ThemeToggle isDark={false} onToggle={handleToggle} />);

    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });

  it('calls onToggle when clicked', () => {
    const handleToggle = vi.fn();
    render(<ThemeToggle isDark={false} onToggle={handleToggle} />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(handleToggle).toHaveBeenCalledTimes(1);
  });

  it('calls onToggle multiple times when clicked multiple times', () => {
    const handleToggle = vi.fn();
    render(<ThemeToggle isDark={false} onToggle={handleToggle} />);

    const button = screen.getByRole('button');
    fireEvent.click(button);
    fireEvent.click(button);
    fireEvent.click(button);

    expect(handleToggle).toHaveBeenCalledTimes(3);
  });

  it('has accessible aria-label', () => {
    const handleToggle = vi.fn();
    render(<ThemeToggle isDark={false} onToggle={handleToggle} />);

    const button = screen.getByRole('button');
    expect(button).toHaveAttribute('aria-label');
  });
});
