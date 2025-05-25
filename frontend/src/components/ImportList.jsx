import React, { useEffect, useState } from 'react';
import {Link} from "react-router-dom";
import axios from 'axios';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSquarePlus, faPlus, faBan, faRetweet, faEquals, faQuestion, faCircleQuestion, faFolderPlus } from '@fortawesome/free-solid-svg-icons';

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
                        <th scope="col">Result</th>
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
                            <td>
                                <Link to={`/channels/${importEntry.channel_id}`} className="text-blue-400 hover:underline">
                                    {importEntry.channel_name}
                                </Link>
                            </td>
                            <td>{importEntry.integration}</td>
                            <td>
                                {importEntry.new_jobs > 0 && <span className="me-1" title="new"><FontAwesomeIcon icon={faPlus} /> {importEntry.new_jobs}</span>}
                                {importEntry.updated_jobs > 0 && <span className="me-1" title="updated"><FontAwesomeIcon icon={faRetweet} /> {importEntry.updated_jobs}</span>}
                                {importEntry.published > 0 && <span className="me-1" title="published"><FontAwesomeIcon icon={faSquarePlus} /> {importEntry.published}</span>}
                                {importEntry.missing_jobs > 0 && <span className="me-1" title="missing"><FontAwesomeIcon icon={faQuestion} /> {importEntry.missing_jobs}</span>}
                                {importEntry.missing_published > 0 && <span className="me-1" title="missing published"><FontAwesomeIcon icon={faCircleQuestion} /> {importEntry.missing_jobs}</span>}
                                {importEntry.failed_jobs > 0 && <span className="me-1" title="failed"><FontAwesomeIcon icon={faBan} /> {importEntry.errors}</span>}
                                {importEntry.late_published > 0 && <span className="me-1" title="late published"><FontAwesomeIcon icon={faFolderPlus} /> {importEntry.late_published}</span>}
                                {importEntry.no_change_jobs > 0 && <span className="me-1" title="not changed"><FontAwesomeIcon icon={faEquals} /> {importEntry.no_change_jobs}</span>}
                            </td>
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
