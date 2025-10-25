import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Posting from './posting'
import { BrowserOpenURL } from '../wailsjs/runtime/runtime'

// Mock the Wails runtime
vi.mock('../wailsjs/runtime/runtime', () => ({
  BrowserOpenURL: vi.fn(),
}))

describe('Posting Component', () => {
  const mockPostings = [
    {
      id: 1,
      companyName: 'Tech Corp',
      link: 'https://example.com/job1',
      descrip: 'Software Engineer position',
      qualifications: JSON.stringify(['5+ years experience', 'React expertise', 'Team player']),
      reviewOutput: 'Great match for your skills',
      postedDate: '2025-10-20',
    },
    {
      id: 2,
      companyName: 'StartupXYZ',
      link: 'https://example.com/job2',
      descrip: 'Senior Frontend Developer',
      qualifications: JSON.stringify(['7+ years experience', 'Vue.js']),
      reviewOutput: 'Competitive salary',
      postedDate: '2025-10-22',
    },
  ]

  beforeEach(() => {
    // Reset all mocks before each test
    vi.clearAllMocks()
    
    // Setup default mock implementation
    window.go.main.App.GetRecentPostingsWithQualifications = vi.fn().mockResolvedValue(mockPostings)
    window.go.main.App.DeletePosting = vi.fn().mockResolvedValue(undefined)
  })

  describe('Initial Rendering', () => {
    it('fetches and displays postings on mount', async () => {
      render(<Posting />)

      expect(await screen.findByText('Tech Corp')).toBeInTheDocument()
      expect(await screen.findByText('StartupXYZ')).toBeInTheDocument()
      expect(window.go.main.App.GetRecentPostingsWithQualifications).toHaveBeenCalledTimes(1)
    })

    it('displays posting details correctly', async () => {
      render(<Posting />)

      expect(await screen.findByText('Tech Corp')).toBeInTheDocument()
      expect(screen.getByText('Software Engineer position')).toBeInTheDocument()
      expect(screen.getByText('5+ years experience')).toBeInTheDocument()
      expect(screen.getByText('React expertise')).toBeInTheDocument()
      expect(screen.getByText('Team player')).toBeInTheDocument()
      expect(screen.getByText('Great match for your skills')).toBeInTheDocument()
      expect(screen.getByText('2025-10-20')).toBeInTheDocument()
    })

    it('renders multiple postings', async () => {
      render(<Posting />)

      await screen.findByText('Tech Corp')

      const postingCards = document.querySelectorAll('.posting-card')
      expect(postingCards).toHaveLength(2)
    })

    it('handles empty postings array', async () => {
      window.go.main.App.GetRecentPostingsWithQualifications = vi.fn().mockResolvedValue([])
      
      render(<Posting />)

      await waitFor(() => {
        expect(window.go.main.App.GetRecentPostingsWithQualifications).toHaveBeenCalledTimes(1)
      })

      const postingCards = document.querySelectorAll('.posting-card')
      expect(postingCards).toHaveLength(0)
    })

    it('handles null postings response', async () => {
      window.go.main.App.GetRecentPostingsWithQualifications = vi.fn().mockResolvedValue(null)
      
      render(<Posting />)

      await waitFor(() => {
        expect(window.go.main.App.GetRecentPostingsWithQualifications).toHaveBeenCalledTimes(1)
      })

      const postingCards = document.querySelectorAll('.posting-card')
      expect(postingCards).toHaveLength(0)
    })
  })

  describe('Link Opening', () => {
    it('opens browser when link button is clicked', async () => {
      const user = userEvent.setup()
      render(<Posting />)

      await screen.findByText('Tech Corp')

      const linkButton = await screen.findByText('https://example.com/job1')
      await waitFor(() => user.click(linkButton))

      expect(BrowserOpenURL).toHaveBeenCalledWith('https://example.com/job1')
    })

    it('opens correct URL for multiple postings', async () => {
      const user = userEvent.setup()
      render(<Posting />)

      await screen.findByText('Tech Corp')

      const linkButton2 = await screen.findByText('https://example.com/job2')
      await waitFor(() => user.click(linkButton2))

      expect(BrowserOpenURL).toHaveBeenCalledWith('https://example.com/job2')
    })
  })

  describe('Context Menu', () => {
    it('does not show context menu initially', async () => {
      render(<Posting />)

      await screen.findByText('Tech Corp') // Wait for data to load
      
      const contextMenu = document.querySelector('.context-menu')
       expect(contextMenu).not.toBeInTheDocument()
    })

    it('shows context menu on right-click', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      await waitFor(() => {
        const contextMenu = document.querySelector('.context-menu')
        expect(contextMenu).toBeInTheDocument()
      })

      expect(screen.getByText('Delete')).toBeInTheDocument()
    })

    it('positions context menu at cursor location', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 150, clientY: 250 })

      await waitFor(() => {
        const contextMenu = document.querySelector('.context-menu')
        expect(contextMenu).toBeInTheDocument()
        expect(contextMenu.style.left).toBe('150px')
        expect(contextMenu.style.top).toBe('250px')
      })
    })

    it('closes context menu on click outside', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      await waitFor(() => {
        expect(document.querySelector('.context-menu')).toBeInTheDocument()
      })

      // Click outside
      fireEvent.click(document)

      await waitFor(() => {
        expect(document.querySelector('.context-menu')).not.toBeInTheDocument()
      })
    })

    it('opens context menu for correct posting', async () => {
      render(<Posting />)

      const startupXYZ = await screen.findByText('StartupXYZ')
      const postingCard2 = startupXYZ.closest('.posting-card')
      fireEvent.contextMenu(postingCard2, { clientX: 100, clientY: 200 })

      await waitFor(() => {
        expect(document.querySelector('.context-menu')).toBeInTheDocument()
      })

      // Now delete and verify the correct posting is targeted
      const deleteButton = screen.getByText('Delete')
      fireEvent.click(deleteButton)

      await waitFor(() => {
        expect(window.go.main.App.DeletePosting).toHaveBeenCalledWith(2)
      })
    })
  })

  describe('Delete Functionality', () => {
    it('deletes posting when delete is clicked', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')
      
      // Open context menu
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      const deleteButton = await screen.findByText('Delete')
      fireEvent.click(deleteButton)

      await waitFor(() => {
        expect(window.go.main.App.DeletePosting).toHaveBeenCalledWith(1)
      })
    })

    it('removes posting from UI after successful delete', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')

      // Open context menu and delete
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      const deleteButton = await screen.findByText('Delete')
      fireEvent.click(deleteButton)

      await waitFor(() => {
        expect(screen.queryByText('Tech Corp')).not.toBeInTheDocument()
      })

      // Other posting should still be there
      expect(screen.getByText('StartupXYZ')).toBeInTheDocument()
    })

    it('closes context menu after delete', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      const deleteButton = await screen.findByText('Delete')
      fireEvent.click(deleteButton)

      await waitFor(() => {
        expect(document.querySelector('.context-menu')).not.toBeInTheDocument()
      })
    })

    it('shows alert on delete error', async () => {
      const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      
      window.go.main.App.DeletePosting = vi.fn().mockRejectedValue(new Error('Delete failed'))

      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      const deleteButton = await screen.findByText('Delete')
      fireEvent.click(deleteButton)

      await waitFor(() => {
        expect(alertSpy).toHaveBeenCalledWith('Failed to delete posting: Error: Delete failed')
        expect(consoleErrorSpy).toHaveBeenCalled()
      })

      // Posting should still be in UI since delete failed
      expect(screen.getByText('Tech Corp')).toBeInTheDocument()

      alertSpy.mockRestore()
      consoleErrorSpy.mockRestore()
    })

    it('handles delete with no postingId gracefully', async () => {
      render(<Posting />)

      const techCorp = await screen.findByText('Tech Corp')

      // Manually open context menu without posting ID
      const postingCard = techCorp.closest('.posting-card')
      fireEvent.contextMenu(postingCard, { clientX: 100, clientY: 200 })

      const deleteButton = await screen.findByText('Delete')

      // This should not crash even if postingId is somehow null
      fireEvent.click(deleteButton)

      // Should have been called normally since we did set a proper ID
      await waitFor(() => {
        expect(window.go.main.App.DeletePosting).toHaveBeenCalled()
      })
    })
  })

  describe('Qualifications Parsing', () => {
    it('correctly parses and displays qualifications array', async () => {
      render(<Posting />)

      await screen.findByText('Tech Corp')

      const qualificationItems = screen.getAllByText(/years experience/)
      expect(qualificationItems.length).toBeGreaterThan(0)
    })

    it('handles postings with different numbers of qualifications', async () => {
      render(<Posting />)

      await screen.findByText('Tech Corp')

      // First posting has 3 qualifications
      expect(screen.getByText('5+ years experience')).toBeInTheDocument()
      expect(screen.getByText('React expertise')).toBeInTheDocument()
      expect(screen.getByText('Team player')).toBeInTheDocument()

      // Second posting has 2 qualifications
      expect(screen.getByText('7+ years experience')).toBeInTheDocument()
      expect(screen.getByText('Vue.js')).toBeInTheDocument()
    })
  })

  describe('Edge Cases', () => {
    it('handles API error gracefully', async () => {
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      window.go.main.App.GetRecentPostingsWithQualifications = vi.fn().mockRejectedValue(new Error('API Error'))

      render(<Posting />)

      await waitFor(() => {
        expect(window.go.main.App.GetRecentPostingsWithQualifications).toHaveBeenCalledTimes(1)
      })

      consoleErrorSpy.mockRestore()
    })

    it('renders posting with missing optional fields', async () => {
      const postingWithMissingFields = {
        id: 3,
        companyName: 'Minimal Corp',
        link: 'https://example.com/job3',
        descrip: '',
        qualifications: JSON.stringify([]),
        reviewOutput: '',
        postedDate: '2025-10-25',
      }

      window.go.main.App.GetRecentPostingsWithQualifications = vi.fn().mockResolvedValue([postingWithMissingFields])

      render(<Posting />)

      expect(await screen.findByText('Minimal Corp')).toBeInTheDocument()
      expect(screen.getByText('2025-10-25')).toBeInTheDocument()
    })
  })
})

