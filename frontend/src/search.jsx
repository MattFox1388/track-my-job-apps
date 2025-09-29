import { useState, useEffect } from 'react'
import './search.css'
import EditModal from './editModal'

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
    const [totalCount, setTotalCount] = useState(0)
    const [editingJobApp, setEditingJobApp] = useState(null)

    useEffect(() => {
        const fetchResults = async () => {
            const results = await window.go.main.App.GetAllJobApps()
            console.log(results)
            // get last 20 results
            setResults(results)
            
            // Get total count from database
            const count = await window.go.main.App.GetJobAppCount()
            setTotalCount(count)
        }
        fetchResults()
    }, [])

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

    const handleResultClick = (app) => {
        console.log(`Editing job app: ${app}`)
        setEditingJobApp(app)
    }

    const handleSaveJob = async (app) => {
        console.log(`Saving job app: ${app}`)
        const updatedApp = await window.go.main.App.UpdateJobApp(app)
        console.log(`Updated job app: ${updatedApp}`)
        setEditingJobApp(null)
    }

    return (
        isLoading ? (
            <div>Loading...</div>
        ) : (
            <>
                <div className="search-container">
                    <div className="search-bar">
                        <div className="search-input-group">
                            <input
                                type="text"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                onKeyDown={handleKeyPress}
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

                    <div className="jobs-summary">
                        <h2>Total Jobs Applied: {totalCount}</h2>
                        {results.length < totalCount && (
                            <p style={{margin: '8px 0 0 0', color: '#666', fontSize: '14px'}}>
                                Showing {results.length} recent applications
                            </p>
                        )}
                    </div>

                    <div className="search-results">
                    {results.map((result) => (
                        <div className="result-item" key={result.appId} onClick={() => handleResultClick(result)}>
                            <h3>{result.company}</h3>
                            <p>{result.position}</p>
                            <p>{result.location}</p>
                            <p>{result.dateApplied}</p>
                            <p>{result.status}</p>
                            <p>{result.notes}</p>
                            <p>{result.website}</p>
                            <p>{result.salaryRange}</p>
                            <p>{result.workplaceType}</p>
                        </div>
                    ))}
                </div>
                </div>
                <EditModal open={editingJobApp !== null} handleClose={() => setEditingJobApp(null)} jobApp={editingJobApp} onSave={handleSaveJob}/>
            </>
        ))
}

export default Search