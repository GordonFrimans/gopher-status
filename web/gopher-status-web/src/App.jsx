import { useState, useEffect } from 'react';
import { Container, Title, Button, Group, ActionIcon, Tooltip } from '@mantine/core';
import { IconPlus, IconLogout } from '@tabler/icons-react';
import { MonitorTable } from './components/MonitorTable';
import { ThemeToggle } from './components/ThemeToggle';
import { UpdateButton } from './components/UpdateButton';
import { AddButton } from './components/AddButton';
// Добавляем импорт страницы авторизации
import AuthPage from './pages/AuthPage';
import { ListMonitors, DeleteMonitor } from './api/api.js';

export default function App() {
    // Добавляем состояние для проверки, залогинен ли юзер
    const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('accessToken'));

    const [monitors, setMonitors] = useState([]);
    const [loading, setLoading] = useState(false);
    const [opened, setOpened] = useState(false);

    // Функция для выхода (Logout)
    const handleLogout = () => {
        localStorage.removeItem('accessToken');
        setIsAuthenticated(false);
    };

    const fetchMonitors = async () => {
        if (!isAuthenticated) return; // Не делаем запросы, если не залогинены

        setLoading(true);
        try {
            const data = await ListMonitors();
            setMonitors(data.monitors || []);
        } catch (error) {
            console.error("Failed to load monitors:", error);
            // Если поймали ошибку авторизации (например, токен истек)
            if (error.message === "Unauthorized") {
                handleLogout();
            }
        } finally {
            setLoading(false);
        }
    };

    // Оставим только один useEffect для чистоты кода
    useEffect(() => {
        if (isAuthenticated) {
            fetchMonitors();
            const intervalId = setInterval(fetchMonitors, 5000);
            return () => clearInterval(intervalId);
        }
    }, [isAuthenticated, opened]); // Перезапустится, если изменится статус логина или модалка

    const handleDelete = async (id) => {
        try {
            await DeleteMonitor(id);
            setMonitors(monitors.filter(m => m.id !== id));
        } catch (error) {
            if (error.message === "Unauthorized") {
                handleLogout();
            }
        }
    };

    // Если нет токена — показываем ТОЛЬКО страницу входа
    if (!isAuthenticated) {
        // Передаем пропс, чтобы AuthPage мог обновить состояние в App.jsx после успешного входа
        return <AuthPage onLoginSuccess={() => setIsAuthenticated(true)} />;
    }

    // Если токен есть — показываем дашборд
    return (
        <Container size="lg" py="xl">
            <Group mb="lg" justify="space-between">
                <Group>
                    <ThemeToggle />
                    <UpdateButton onUpdate={fetchMonitors} loading={loading} />
                </Group>
                {/* Кнопка выхода */}
                <Tooltip label="Выйти из аккаунта" position="bottom" withArrow>
                    <ActionIcon
                        color="red"
                        variant="light"
                        onClick={handleLogout}
                        size="lg"
                        radius="xl" // Делает кнопку идеально круглой
                    >
                        <IconLogout stroke={1.5} />
                    </ActionIcon>
                </Tooltip>
            </Group>

            <Group justify="space-between" mb="lg">
                <Title order={2}>GopherStatus Dashboard 🐹</Title>
                <AddButton opened={opened} setOpened={setOpened} />
            </Group>

            <MonitorTable monitors={monitors} onDelete={handleDelete} />
        </Container>
    );
}
