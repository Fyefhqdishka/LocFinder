import { useEffect, useState } from "react";

const Main = () => {
    const [address, setAddress] = useState("");
    const [location, setLocation] = useState(null);
    const [locations, setLocations] = useState([]); // Состояние для всех локаций

    // Получение одной локации по IP
    const fetchLocation = async (address) => {
        try {
            const response = await fetch(`
                http://localhost:8000/location?ip=${encodeURIComponent(address)}`,
            {
                method: "GET",
                    headers: {
                "Content-Type": "application/json",
            },
                mode: "cors",
            }
        );
            if (!response.ok) {
                throw new Error("Failed to fetch location");
            }
            const data = await response.json();
            return data;
        } catch (error) {
            console.error(error);
            return null;
        }
    };

    // Получение всех локаций
    const fetchLocations = async () => {
        try {
            const response = await fetch(`http://localhost:8000/locations`, {
            method: "GET",
                headers: {
                "Content-Type": "application/json",
            },
            mode: "cors",
        });
        if (!response.ok) {
            throw new Error("Failed to fetch locations");
        }
        const data = await response.json();
        setLocations(data.result); // Сохраняем список локаций в состояние
    } catch (error) {
        console.error(error);
    }
};

const handleSearch = async () => {
    if (!address.trim()) return;
    const locationData = await fetchLocation(address);
    setLocation(locationData?.result || null); // Сохраняем результат только из "result"
};


const fetchDefault = async () => {
    try {
        const response = await fetch(
            `http://localhost:8000/location`,
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
                mode: "cors",
            }
        );
        if (!response.ok) {
            throw new Error("Failed to fetch location");
        }
        const data = await response.json();
        return data;
    } catch (error) {
        console.error(error);
        return null;
    }
};

const handleDefault = async () => {
    if (!address.trim()) return;
    const locationData = await fetchDefault();
    setLocation(locationData?.result || null); // Сохраняем результат только из "result"
};

useEffect(()=> {
    handleDefault();
}, [location])

return (
    <div>
        <input
            type="text"
            value={address}
            onChange={(e) => setAddress(e.target.value)}

        />
        <button onClick={handleSearch}>Search</button>
        {location && (
            <div>
                <pre>{JSON.stringify(location, null, 2)}</pre>
            </div>
        )}

        <button onClick={fetchLocations}>Show All Locations</button>
        <div>
            {locations.length > 0 && (
                <ul>
                    {locations.map((loc, index) => (
                        <li key={index}>
                            <pre>{JSON.stringify(loc, null, 2)}</pre>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    </div>
);
};

export default Main;