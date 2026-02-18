import { IconPlus } from '@tabler/icons-react';
import { useState } from 'react';
import { Modal, Button, TextInput, NumberInput } from '@mantine/core';
import { CreateMonitor } from '../api/api.js';
export function AddButton({ opened, setOpened }) {
    const [formData, setFormData] = useState({
        name: '',
        url: '',
        interval: 30
    });
    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            // Тут вызов твоего API для создания монитора
            console.log('Creating monitor:', formData);
            await CreateMonitor(formData);

            // После успешного создания:
            setOpened(false); // Закрыть модалку
            setFormData({ name: '', url: '', interval: 30 }); // Очистить форму
        } catch (error) {
            console.error('Failed to create monitor:', error);
        }
    };


    return (
        <>
            <Modal
                opened={opened}
                onClose={() => setOpened(false)}
                title="Add New Monitor 🐹"
                centered // Центрирует по вертикали!
            >
                <form onSubmit={handleSubmit}>
                    <TextInput
                        label="Monitor Name"
                        placeholder="Enter name"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}

                        required
                    />
                    <TextInput
                        label="URL"
                        placeholder="https://example.com"
                        value={formData.url}
                        onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                        required
                        mt="md"
                    />
                    <NumberInput
                        label="Check Interval (seconds)"
                        placeholder="30"
                        value={formData.interval}
                        onChange={(value) => setFormData({ ...formData, interval: value })}
                        min={10}
                        max={3600}
                        mt="md"
                    />
                    <Button type="submit" fullWidth mt="xl">
                        Create Monitor
                    </Button>
                </form>
            </Modal>


            <Button
                leftSection={<IconPlus size={16} />}
                variant="filled"
                onClick={() => setOpened(true)}
            >
                Add Monitor
            </Button>
        </>
    );
}

