import { useEffect, useState } from "react";
import styles from "../page/Main.module.css";
import { AiFillCaretDown } from "react-icons/ai";

const Main = () => {
    const [address, setAddress] = useState("");
    const [location, setLocation] = useState(null);
    const [locations, setLocations] = useState([]); // Состояние для всех локаций
    const [isOpen, setIsOpen] = useState(false);

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


    const fetchLocations = async () => {
        try {
            const response = await fetch('http://localhost:8000/locations', {
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

    const fetchDefault = async () => {
        try {
            const response = await fetch(
                "http://localhost:8000/location",
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
            setAddress(data.result.query || "")
            return data
        } catch (error) {
            console.error(error);
            return null;
        }
    };

    const handleSearch = async () => {
        if (!address.trim()) return;
        const locationData = await fetchLocation(address);
        setLocation(locationData?.result || null); // Сохраняем результат только из "result"
    };

    const handleDefault = async () => {
        const locationData = await fetchDefault();
        setLocation(locationData?.result || null); // Сохраняем результат только из "result"
    };

    const handleKeyPress = (e) => {
        if (e.key === "Enter") { // Проверяем, нажата ли клавиша Enter
            handleSearch();
        }
    };

    useEffect(()=> {
        handleDefault();
    }, [])

    return (
        <div className={styles.mainPage}>
            <section>
                <input
                    type="text"
                    value={address}
                    onChange={(e) => setAddress(e.target.value)}
                    onKeyDown={handleKeyPress}
                />
                <button className={styles.searchBtn} onClick={handleSearch}>Search</button>
                {location && (
                    <div className={styles.location}>
                        <pre>{JSON.stringify(location, null, 2)}</pre>
                    </div>
                )}
            </section>

            <section className={`${isOpen ? styles.show : styles.nope}`}>
                <div className={styles.buttons}>
                    <button onClick={fetchLocations}>Show All Locations</button>
                    <button onClick={()=> setIsOpen(!isOpen)}>
                        <AiFillCaretDown />
                    </button>
                </div>
                <div className={styles.locations}>
                    {locations.length > 0 && (
                        <div className={styles.location}>
                            <ul>
                                {locations.map((loc, index) => (
                                    <li key={index}>
                                        <pre>{JSON.stringify(loc, null, 2)}</pre>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            </section>

        </div>
    );
};

export default Main;