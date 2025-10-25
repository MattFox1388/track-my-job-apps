import '@testing-library/jest-dom'
import { expect, afterEach, vi } from 'vitest'
import { cleanup } from '@testing-library/react'

// Cleanup after each test
afterEach(() => {
  cleanup()
})

// Mock Wails runtime
global.window = global.window || {}
window.go = {
  main: {
    App: {
      GetRecentPostingsWithQualifications: vi.fn(),
      DeletePosting: vi.fn(),
    }
  }
}

