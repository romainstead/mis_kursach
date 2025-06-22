import React, { useEffect, useState } from 'react';
import api from "../utils/api";
import './Dashboard.css';
import { Doughnut } from 'react-chartjs-2';
import {
    Chart as ChartJS,
    ArcElement,
    Tooltip,
    Legend
} from 'chart.js';

ChartJS.register(ArcElement, Tooltip, Legend);


function Dashboard() {
    const [metrics, setMetrics] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        api.get('http://127.0.0.1:8080/api/SetMetrics')
            .then(res => {
                const data = res.data;
                setMetrics({
                    occupancy: data.occupancy,
                    unpaidBookings: data.unpaid_bookings,
                    currentBookings: data.current_bookings,
                    openComplaints: data.open_complaints,
                    freeRooms: data.free_rooms,
                    maintenanceRooms: data.rooms_under_maintenance,
                    revenueLast7Days: data.revenue_7_days,
                    revPar: data.revpar,
                    newGuests7Days: data.new_guests_7_days,
                    revPac: data.revpac,
                });
            })
            .catch(err => {
                console.error('Ошибка при запросе:', err);
            });
    }, []);


    if (!metrics) {
        return (
            <div className="dashboard loader-container">
                <div className="loader"></div>
            </div>
        );
    }
    const doughnutData = {
        datasets: [
            {
                data: [metrics.occupancy, 100 - metrics.occupancy],
                backgroundColor: ['#FBBA00', '#e0e0e0'],
                borderWidth: 0,
            },
        ],
    };
    return (
        <div className="dashboard">
            <div className="big-cards">
                <div className="card big">
                    <h3>Текущая загрузка</h3>
                    <div className="chart-wrapper">
                        <Doughnut data={doughnutData} options={{cutout: '70%'}}/>
                        <div className="chart-center">{metrics.occupancy}%</div>
                    </div>
                </div>
                <div className="card big">
                    <h3>Неоплаченные бронирования</h3>
                    <p>{metrics.unpaidBookings}</p>
                </div>
                <div className="card big">
                    <h3>Текущие бронирования</h3>
                    <p>{metrics.currentBookings}</p>
                </div>
                <div className="card big">
                    <h3>Открытые жалобы</h3>
                    <p>{metrics.openComplaints}</p>
                </div>
            </div>
            <div className="small-cards">
                <div className="card small alter"><h4>Свободные номера</h4><p>{metrics.freeRooms}</p></div>
                <div className="card small"><h4>В обслуживании</h4><p>{metrics.maintenanceRooms}</p></div>
                <div className="card small alter"><h4>Выручка (7 дней)</h4>
                    <p>{metrics.revenueLast7Days.toLocaleString()}₽</p></div>
                <div className="card small"><h4>RevPAR</h4><p>{metrics.revPar}₽</p></div>
                <div className="card small alter"><h4>Новые гости (7 дней)</h4><p>{metrics.newGuests7Days}</p></div>
                <div className="card small"><h4>RevPAC</h4><p>{metrics.revPac}₽</p></div>
            </div>
        </div>
    );
}

export default Dashboard;
