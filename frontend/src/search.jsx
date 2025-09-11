import { useState } from 'react'
import './search.css'

const SearchType = Object.freeze({
    COMPANY: 'company',
    POSITION: 'position',
    FULL_TEXT: 'full text',
})

function Search() {
    const [searchTerm, setSearchTerm] = useState('')
    const [searchType, setSearchType] = useState(SearchType.COMPANY)
    const [results, setResults] = useState([])
    const [isLoading, setIsLoading] = useState(false)

    const handleSearch = async () => {
        setIsLoading(true)
        try {
            const results = await window.go.main.App.SearchByCompany(searchTerm)
            setResults(results)
        } catch (error) {
            console.error("Error searching:", error)
        } finally {
            setIsLoading(false)
        }
    }

    const handleKeyPress = (e) => {
        if (e.key === 'Enter') {
            handleSearch()
        }
    }

    return (
        isLoading ? (
            <div>Loading...</div>
        ) : (
            <div className="search-container">
                <div className="search-bar">
                    <div className="search-input-group">
                        <input
                            type="text"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            onKeyPress={handleKeyPress}
                            placeholder="Enter company name..."
                            className="search-input"
                            disabled={isLoading}
                        />
                        <button
                            onClick={handleSearch}
                            disabled={isLoading || !searchTerm.trim()}
                            className="search-button"
                        >
                            {isLoading ? '⏳' : '🔍'}
                        </button>
                    </div>

                    <div className="search-type-selector">
                        <label>Search by:</label>
                        <select
                            value={searchType}
                            onChange={(e) => setSearchType(e.target.value)}
                            className="search-select"
                        >
                            <option value={SearchType.COMPANY}>Company</option>
                            <option value={SearchType.POSITION}>Position</option>
                            <option value={SearchType.FULL_TEXT}>Full Text</option>
                        </select>
                    </div>
                </div>
            </div>
        ))
}

export default Search