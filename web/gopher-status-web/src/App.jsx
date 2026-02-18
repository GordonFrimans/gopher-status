import { useState, useEffect } from 'react';
import { Container, Title, Button, Group } from '@mantine/core';
import { MonitorTable } from './components/MonitorTable';
import { ThemeToggle } from './components/ThemeToggle';
import { UpdateButton } from './components/UpdateButton';
import { AddButton } from './components/AddButton';
import { IconPlus } from '@tabler/icons-react';
import { ListMonitors } from './api/api.js';
import { DeleteMonitor } from './api/api.js';

export default function App() {
    const [monitors, setMonitors] = useState([]);
    const [loading, setLoading] = useState(false);
    const [opened, setOpened] = useState(false);


    // Функция для загрузки данных
    const fetchMonitors = async () => {
        setLoading(true);
        try {
            const data = await ListMonitors();
            setMonitors(data.monitors || []); // Адаптируй под структуру твоего ответа
            console.log("Monitors loaded:", data);
        } catch (error) {
            console.error("Failed to load monitors:", error);
        } finally {
            setLoading(false);
        }
    };

    // Загрузить данные при первом рендере

    useEffect(() => {
        fetchMonitors();
    }, []);
    useEffect(() => {
        fetchMonitors();
    }, [opened]);

    useEffect(() => {
        fetchMonitors();

        const intervalId = setInterval(fetchMonitors, 5000);

        return () => clearInterval(intervalId);
    }, []);

    const handleDelete = (id) => {
        console.log("Deleting monitor:", id);
        DeleteMonitor(id)
        setMonitors(monitors.filter(m => m.id !== id));
    };

    return (
        <Container size="lg" py="xl">
            <Group mb="lg">
                <ThemeToggle />
                <UpdateButton onUpdate={fetchMonitors} loading={loading} />
            </Group>

            <Group justify="space-between" mb="lg">
                <Title order={2}>GopherStatus Dashboard 🐹</Title>

                <AddButton opened={opened} setOpened={setOpened} />
            </Group>

            <MonitorTable monitors={monitors} onDelete={handleDelete} />
        </Container>
    );
}
