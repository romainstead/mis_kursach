import {BrowserRouter as Router, Routes, Route, useLocation, Navigate} from 'react-router-dom';
import { useEffect } from 'react';
import Bookings from './components/Bookings';
import Complaints from './components/Complaints';
import Payments from './components/Payments';
import Navbar from './components/Navbar';
import Rooms from './components/Rooms';
import Main from './components/Main';
import Footer from './components/Footer';
import Login from './components/Login';
import Logout from "./components/Logout";
import PrivateRoute from "./components/PrivateRoute";
import CreateBookingForm from "./components/CreateBooking";
import BookingDetails from "./components/BookingDetails";
import ComplaintDetails from "./components/ComplaintDetails";
import CreateComplaintForm from "./components/CreateComplaint";
import {NotificationProvider} from "./components/NotificationContext";

function Layout() {
    const location = useLocation();
    const hideLayout = ["/login", "/logout"].includes(location.pathname);

    useEffect(() => {
        window.scrollTo(0, 0);
    }, [location.pathname]);

    return (
        <NotificationProvider>
            <>
                {!hideLayout && <Navbar />}
                <div className="App" style={{ minHeight: 'calc(100vh - 100px)' }}>
                    <Routes>
                        <Route
                            path="/"
                            element={
                                <Navigate to="/login" />
                            }
                        />
                        <Route path="/main" element={<PrivateRoute><Main /></PrivateRoute>} />
                        <Route path="/bookings" element={<PrivateRoute><Bookings /></PrivateRoute>} />
                        <Route path="/complaints" element={<PrivateRoute><Complaints /></PrivateRoute>} />
                        <Route path="/payments" element={<PrivateRoute><Payments /></PrivateRoute>} />
                        <Route path="/rooms" element={<PrivateRoute><Rooms /></PrivateRoute>} />
                        <Route path="/login" element={<Login />} />
                        <Route path="/logout" element={<Logout />} />
                        <Route path="/create-booking" element={<CreateBookingForm />} />
                        <Route path="/bookings/:id" element={<BookingDetails />} />
                        <Route path="/complaints/:id" element={<ComplaintDetails/>} />
                        <Route
                            path="/create-complaint"
                            element={<CreateComplaintForm onClose={() => {}} />}
                        />
                    </Routes>
                </div>
                {!hideLayout && <Footer />}
            </>
        </NotificationProvider>
    );
}

function App() {
    return (
        <Router>
            <Layout />
        </Router>
    );
}

export default App;
