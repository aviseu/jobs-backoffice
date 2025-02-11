import React, { useEffect, useState } from 'react';
import {Link} from "react-router-dom";
import axios from 'axios';

const ChannelList = () => {
    const [channels, setChannels] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchChannels = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/channels`);
                const data = response.data.channels;
                if (Array.isArray(data)) {
                    setChannels(data);
                } else {
                    setError('Unexpected non array response')
                    console.error('Expected an array but got:', data);
                }
            } catch (err) {
                console.log(err);
                if (err.response) {
                    setError(err.response.data.error.message || "Submission failed. Please check your input.");
                } else if (err.request) {
                    setError("Network error. Please check your connection.");
                } else {
                    setError("An unexpected error occurred.");
                }
            } finally {
                setLoading(false);
            }
        }

        fetchChannels();
    }, []);

    if (loading) return <div className="text-center mt-5"><div className="spinner-border" role="status"></div></div>;

    return (
        <div className="table-responsive mt-5">
            {error && <div className="alert alert-danger" role="alert">
                {error}
            </div>}
            <table className="table table-striped">
                <thead>
                <tr>
                    <th scope="col">Name</th>
                    <th scope="col">Integration</th>
                    <th scope="col">Status</th>
                </tr>
                </thead>
                <tbody>
                {channels.map((channel) => (
                    <tr key={channel.id} >
                        <td>
                            <Link to={`/${channel.id}`} className="text-blue-400 hover:underline">
                                {channel.name}
                            </Link>
                        </td>
                        <td >{channel.integration}</td>
                        <td>{channel.status}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default ChannelList;
