// src/App.js
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './components/Home';
import Bookings from './components/Bookings';
import Complaints from "./components/Complaints";
import Payments from "./components/Payments";
import Navbar from './components/Navbar';

function App() {
    return (
        <Router>
            <Navbar/>
            <div className="App">
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/bookings" element={<Bookings />} />
                    <Route path="/complaints" element={<Complaints />} />
                    <Route path="/payments" element={<Payments/>}/>
                </Routes>
            </div>
        </Router>
    );
}

export default App;