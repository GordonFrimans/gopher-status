import { useState } from 'react';
import { RegisterUser, LoginUser } from '../api/api.js';
import {
    Container,
    Paper,
    Title,
    TextInput,
    PasswordInput,
    Button,
    Text,
    Anchor,
    Stack,
    Alert
} from '@mantine/core';

export default function AuthPage({ onLoginSuccess }) {
    const [isLoginMode, setIsLoginMode] = useState(true); // true = логин, false = регистрация
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const [successMsg, setSuccessMsg] = useState(null);
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setSuccessMsg(null);
        setLoading(true);

        try {
            if (isLoginMode) {
                // УСПЕШНЫЙ ВХОД
                const data = await LoginUser(login, password);
                localStorage.setItem('accessToken', data.jwt);
                onLoginSuccess();
            } else {
                // УСПЕШНАЯ РЕГИСТРАЦИЯ
                await RegisterUser(login, password);
                setSuccessMsg('Пользователь успешно создан! Теперь войдите.');
                setIsLoginMode(true);
                setPassword('');
            }
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Container size={420} my={80}>
            <Title ta="center" order={2}>
                {isLoginMode ? 'Вход в GopherStatus 🐹' : 'Регистрация'}
            </Title>

            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                {/* Блок для вывода ошибок */}
                {error && (
                    <Alert color="red" mb="md" title="Ошибка">
                        {error}
                    </Alert>
                )}

                {/* Блок для вывода успешной регистрации */}
                {successMsg && (
                    <Alert color="green" mb="md" title="Успех">
                        {successMsg}
                    </Alert>
                )}

                <form onSubmit={handleSubmit}>
                    <Stack>
                        <TextInput
                            label="Логин"
                            placeholder="Введите ваш логин"
                            required
                            value={login}
                            onChange={(e) => setLogin(e.target.value)}
                        />
                        <PasswordInput
                            label="Пароль"
                            placeholder="Введите ваш пароль"
                            required
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        <Button type="submit" fullWidth mt="xl" loading={loading}>
                            {isLoginMode ? 'Войти' : 'Создать аккаунт'}
                        </Button>
                    </Stack>
                </form>

                <Text c="dimmed" size="sm" ta="center" mt={20}>
                    {isLoginMode ? 'Нет аккаунта? ' : 'Уже есть аккаунт? '}
                    <Anchor
                        size="sm"
                        component="button"
                        type="button"
                        onClick={() => {
                            setIsLoginMode(!isLoginMode);
                            setError(null);
                            setSuccessMsg(null);
                        }}
                    >
                        {isLoginMode ? 'Зарегистрироваться' : 'Войти'}
                    </Anchor>
                </Text>
            </Paper>
        </Container>
    );
}
