import { useEffect } from "react";
import api from "../utils/api";

function Logout() {
    useEffect(() => {
        api.post('http://127.0.0.1:8080/api/logout')
            .finally(() => {
                localStorage.removeItem("token");
                window.location.href = "/login";
            });
    }, []);

    return null;
}

export default Logout;
