import '@mantine/core/styles.css'; // CRITICAL
import { MantineProvider } from '@mantine/core';
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.jsx';

ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        <MantineProvider>
            <App />
        </MantineProvider>
    </React.StrictMode>,
);
