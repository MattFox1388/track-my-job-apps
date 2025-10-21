import { useState, useEffect } from 'react'
import './posting.css'

// Dummy data
const DUMMY_POSTINGS = [
    {
        id: 1,
        link: 'https://example.com/job/1',
        descrip: 'We are seeking a talented Software Engineer to join our growing team. You will work on cutting-edge technologies and contribute to building scalable applications.',
        qualifications: ['5+ years experience in software development', 'Proficiency in React and Node.js', 'Strong problem-solving skills', 'Bachelor\'s degree in Computer Science'],
        postedDate: '2025-10-01',
        companyName: 'TechCorp Solutions',
        reviewOutput: 'Strong match - Your experience aligns well with the requirements. The company culture emphasizes innovation and work-life balance.'
    },
    {
        id: 2,
        link: 'https://example.com/job/2',
        descrip: 'Looking for a Senior Backend Developer to architect and build robust API services. You\'ll be working with microservices and cloud infrastructure.',
        qualifications: ['7+ years backend development', 'Experience with Go or Python', 'Knowledge of AWS/GCP', 'Database design expertise'],
        postedDate: '2025-09-28',
        companyName: 'CloudFirst Inc',
        reviewOutput: 'Good opportunity - Competitive salary range. Remote-friendly position with flexible hours.'
    },
    {
        id: 3,
        link: 'https://example.com/job/3',
        descrip: 'Join our startup as a Full Stack Developer! Help us build the next generation of fintech solutions.',
        qualifications: ['3+ years full stack development', 'React and TypeScript', 'REST API design', 'Agile methodology experience'],
        postedDate: '2025-10-05',
        companyName: 'FinTech Innovations',
        reviewOutput: 'Moderate match - Startup environment with high growth potential but may require longer hours.'
    }
]

function Posting() {
    const [postings, setPostings] = useState([])
    const [contextMenu, setContextMenu] = useState({ visible: false, x: 0, y: 0, postingId: null })

    useEffect(() => {
        const fetchPostings = async () => { 
            const results = await window.go.main.App.GetRecentPostingsWithQualifications();
            results.map(result => {
                result.qualifications = JSON.parse(result.qualifications)
            })
        
            setPostings(results || [])
        }
        fetchPostings()
    }, [])

    useEffect(() => {
        // Close context menu when clicking anywhere
        const handleClick = () => setContextMenu(prev => ({ ...prev, visible: false }))
        document.addEventListener('click', handleClick)
        return () => document.removeEventListener('click', handleClick)
    }, [])

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
        <div className="posting-container">
            <h1>Job Postings</h1>
            <div className="postings-list">
                {postings.map((posting) => (
                    <div 
                        key={posting.id} 
                        className="posting-card"
                        onContextMenu={(e) => handleContextMenu(e, posting.id)}
                    >
                        <h2 className="posting-title">{posting.companyName}</h2>
                        <a 
                            href={posting.link} 
                            target="_blank" 
                            rel="noopener noreferrer"
                            className="posting-link"
                        >
                            {posting.link}
                        </a>
                        
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