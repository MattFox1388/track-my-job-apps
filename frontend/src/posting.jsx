import { useState, useEffect, useLayoutEffect } from 'react'
import { BrowserOpenURL} from '../wailsjs/runtime/runtime'
import './posting.css'

function Posting() {
    const [postings, setPostings] = useState([])
    const [contextMenu, setContextMenu] = useState({ visible: false, x: 0, y: 0, postingId: null })
    // Initialize as true to prevent saving until after restoration
    const [isRestoring, setIsRestoring] = useState(() => {
        return sessionStorage.getItem('postingsScrollPos') !== null
    })

    useEffect(() => {
        console.log('Posting component mounted...')
        const fetchPostings = async () => { 
            try {
                const results = await window.go.main.App.GetRecentPostingsWithQualifications() ||  [];
                results.map(result => {
                    result.qualifications = typeof result.qualifications === 'string' ? JSON.parse(result.qualifications) : result.qualifications
                })
            
                setPostings(results || [])
            } catch (error) {
                console.error('Failed to fetch postings:', error)
                setPostings([])
            }
           
        }
        fetchPostings()
    }, [])

    useEffect(() => {
        // Close context menu when clicking anywhere
        const handleClick = () => setContextMenu(prev => ({ ...prev, visible: false }))
        document.addEventListener('click', handleClick)
        return () => document.removeEventListener('click', handleClick)
    }, [])

    // Restore scroll position when postings are loaded (useLayoutEffect runs synchronously after DOM updates)
    useLayoutEffect(() => { 
        if (postings.length === 0) return
        
        const savedScrollPos = sessionStorage.getItem('postingsScrollPos')
        console.log('Restoring scroll to:', savedScrollPos)
        
        if (savedScrollPos) {
            const targetScroll = parseInt(savedScrollPos, 10)
            
            // Use setTimeout to ensure DOM is fully painted
            setTimeout(() => {
                window.scrollTo(0, targetScroll)
                console.log('Restored to window.scrollY:', window.scrollY)
                
                // Allow saving again after restoration completes
                setTimeout(() => {
                    console.log('Restoration complete, enabling scroll saving')
                    setIsRestoring(false)
                }, 500)
            }, 0)
        } else {
            // No saved position, enable scroll saving immediately
            setIsRestoring(false)
        }
    }, [postings])

    // Save scroll position continuously as user scrolls (but not during restoration)
    useEffect(() => {
        const saveScroll = () => {
            if (isRestoring) {
                console.log('Skipping save during restoration')
                return
            }
            const pos = window.scrollY
            console.log('Saving scroll position:', pos)
            sessionStorage.setItem('postingsScrollPos', pos)
        }
        
        // Throttle to avoid excessive saves
        let timeoutId
        const throttledSave = () => {
            clearTimeout(timeoutId)
            timeoutId = setTimeout(saveScroll, 150)
        }
        
        window.addEventListener('scroll', throttledSave)
        
        return () => {
            window.removeEventListener('scroll', throttledSave)
            clearTimeout(timeoutId)
        }
    }, [isRestoring])

    const handleContextMenu = (e, postingId) => {
        e.preventDefault()
        setContextMenu({
            visible: true,
            x: e.clientX,
            y: e.clientY,
            postingId: postingId
        })
    }

    const handleDelete = async () => {
        if (!contextMenu.postingId) return
        
        try {
            await window.go.main.App.DeletePosting(contextMenu.postingId)
            // Remove from local state
            setPostings(postings.filter(p => p.id !== contextMenu.postingId))
            setContextMenu({ ...contextMenu, visible: false })
        } catch (error) {
            console.error('Failed to delete posting:', error)
            alert('Failed to delete posting: ' + error)
        }
    }

    return (
        <div className="posting-container" >
            <h1>Job Postings</h1>
            <div className="postings-list">
                {postings.map((posting) => (
                    <div 
                        key={posting.id} 
                        className="posting-card"
                        onContextMenu={(e) => handleContextMenu(e, posting.id)}
                    >
                        <h2 className="posting-title">{posting.companyName}</h2>
                        <button 
                            onClick={() => BrowserOpenURL(posting.link)}
                            className="posting-link posting-link-button"
                        >
                            {posting.link}
                        </button>
                        <div className="posting-body">
                            <p className="posting-description">{posting.descrip}</p>
                            
                            <div className="posting-qualifications">
                                <strong>Qualifications:</strong>
                                <ul>
                                    {posting.qualifications.map((qual, idx) => (
                                        <li key={idx}>{qual}</li>
                                    ))}
                                </ul>
                            </div>
                            
                            <div className="posting-review">
                                <strong>Review:</strong>
                                <p>{posting.reviewOutput}</p>
                            </div>
                        </div>
                        
                        <p className="posting-date">{posting.postedDate}</p>
                    </div>
                ))}
            </div>

            {contextMenu.visible && (
                <div 
                    className="context-menu"
                    style={{
                        position: 'fixed',
                        top: `${contextMenu.y}px`,
                        left: `${contextMenu.x}px`,
                        zIndex: 1000
                    }}
                >
                    <button onClick={handleDelete}>Delete</button>
                </div>
            )}
        </div>
    )
}

export default Posting