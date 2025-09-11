import { Link, Routes, Route } from 'react-router-dom'
import TrackJob from './TrackJob'
import Search from './search'
import './App.css'

function App() {
    return (
        <div className="app-container">
            <div>
                <nav>
                    <Link to="/">Track Job</Link>
                    <Link to="/search">Search</Link>
                </nav>
            </div>


            <div className="main-content">
                <Routes>
                    <Route path="/" index element={<TrackJob />} />
                    <Route path="/search" element={<Search />} />
                </Routes>
            </div>

        </div>
    )
}

export default App