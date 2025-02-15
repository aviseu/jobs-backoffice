import React, { useEffect, useState } from 'react';
import {Link} from "react-router-dom";
import axios from 'axios';

const ImportList = () => {
    const [imports, setImports] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchImports = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/api/imports`);
                const data = response.data.imports;
                if (Array.isArray(data)) {
                    setImports(data);
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

        fetchImports();
    }, []);

    if (loading) return <div className="text-center mt-5"><div className="spinner-border" role="status"></div></div>;

    return (
        <div className="mt-5">
            <div className="table-responsive">
                {error && <div className="alert alert-danger" role="alert">
                    {error}
                </div>}
                <table className="table table-striped">
                    <thead>
                    <tr>
                        <th scope="col">ID</th>
                        <th scope="col">Channel</th>
                        <th scope="col">Integration</th>
                        <th scope="col">Count</th>
                        <th scope="col">Status</th>
                        <th scope="col">Started</th>
                        <th scope="col">Completed</th>
                    </tr>
                    </thead>
                    <tbody>
                    {imports.map((importEntry) => (
                        <tr key={importEntry.id} >
                            <td>
                                <Link to={`/imports/${importEntry.id}`} className="text-blue-400 hover:underline">
                                    {importEntry.id}
                                </Link>
                            </td>
                            <td>-</td>
                            <td>-</td>
                            <td>{importEntry.total_jobs}</td>
                            <td>{importEntry.status}</td>
                            <td>{importEntry.started_at}</td>
                            <td>{importEntry.ended_at}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default ImportList;
