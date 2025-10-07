import { useState, useEffect } from 'react'
import './posting.css'

// Dummy data
const DUMMY_POSTINGS = [
    {
        id: 1,
        link: 'https://example.com/job/1',
        descrip: 'We are seeking a talented Software Engineer to join our growing team. You will work on cutting-edge technologies and contribute to building scalable applications.',
        qualifs: ['5+ years experience in software development', 'Proficiency in React and Node.js', 'Strong problem-solving skills', 'Bachelor\'s degree in Computer Science'],
        postedDate: '2025-10-01',
        companyName: 'TechCorp Solutions',
        reviewOutput: 'Strong match - Your experience aligns well with the requirements. The company culture emphasizes innovation and work-life balance.'
    },
    {
        id: 2,
        link: 'https://example.com/job/2',
        descrip: 'Looking for a Senior Backend Developer to architect and build robust API services. You\'ll be working with microservices and cloud infrastructure.',
        qualifs: ['7+ years backend development', 'Experience with Go or Python', 'Knowledge of AWS/GCP', 'Database design expertise'],
        postedDate: '2025-09-28',
        companyName: 'CloudFirst Inc',
        reviewOutput: 'Good opportunity - Competitive salary range. Remote-friendly position with flexible hours.'
    },
    {
        id: 3,
        link: 'https://example.com/job/3',
        descrip: 'Join our startup as a Full Stack Developer! Help us build the next generation of fintech solutions.',
        qualifs: ['3+ years full stack development', 'React and TypeScript', 'REST API design', 'Agile methodology experience'],
        postedDate: '2025-10-05',
        companyName: 'FinTech Innovations',
        reviewOutput: 'Moderate match - Startup environment with high growth potential but may require longer hours.'
    }
]

function Posting() {
    const [postings, setPostings] = useState([])

    useEffect(() => {

        const fetchPostings = async () => { 
            await window.go.main.App.GetRecentPostingsWithQualifications();
        
            // Using dummy data for now
            setPostings(DUMMY_POSTINGS)
        }
        // TODO: Replace with actual Wails backend call
        // const fetchPostings = async () => {
        //     const data = await window.go.main.App.GetJobPostings()
        //     setPostings(data || [])
        // }
        // fetchPostings()
        fetchPostings()
    }, [])

    return (
        <div className="posting-container">
            <h1>Job Postings</h1>
            <div className="postings-list">
                {postings.map((posting) => (
                    <div key={posting.id} className="posting-card">
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
                                    {posting.qualifs.map((qual, idx) => (
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
        </div>
    )
}

export default Posting