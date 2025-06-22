// src/App.js
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Bookings from './components/Bookings';
import Complaints from "./components/Complaints";
import Payments from "./components/Payments";
import Navbar from './components/Navbar';
import Rooms from "./components/Rooms";
import Main from "./components/Main";
import Footer from "./components/Footer";

function App() {
    return (
        <Router>
            <Navbar/>
            <div className="App">
                <Routes>
                    <Route path="/" element={<Main />} />
                    <Route path="/bookings" element={<Bookings />} />
                    <Route path="/complaints" element={<Complaints />} />
                    <Route path="/payments" element={<Payments/>}/>
                    <Route path="/rooms" element={<Rooms />} />
                </Routes>
            </div>
            <Footer/>
        </Router>
    );
}

export default App;