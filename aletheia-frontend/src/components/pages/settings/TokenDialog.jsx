import React, { useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useSnackbar } from 'notistack';

function TokenDialog({ token, originalToken, expirationDate, onClose }) {
    const { enqueueSnackbar } = useSnackbar();

    const copyToClipboard = () => {
        // Копируем исходный токен с точками
        navigator.clipboard.writeText(originalToken);
        enqueueSnackbar('Токен скопирован в буфер обмена', { variant: 'success' });
    };

    useEffect(() => {
        const handleEsc = (event) => {
            if (event.key === 'Escape') {
                onClose();
            }
        };
        window.addEventListener('keydown', handleEsc);
        return () => window.removeEventListener('keydown', handleEsc);
    }, [onClose]);

    const modalContent = (
        <div
            className="fixed inset-0 flex items-center justify-center bg-black/50 z-50"
            onClick={onClose}
        >
            <div
                className="bg-white p-6 rounded-lg shadow-lg"
                onClick={(e) => e.stopPropagation()}
            >
                <h2 className="text-xl font-bold mb-4">Ваш токен</h2>
                <p className="mb-2">
                    Токен:
                    <pre className="mt-2 p-2 bg-gray-100 rounded whitespace-pre-wrap">{token}</pre>
                </p>
                <p className="mb-4">Действителен до: {expirationDate}</p>
                <button
                    onClick={copyToClipboard}
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                    Копировать токен
                </button>
                <button
                    onClick={onClose}
                    className="ml-4 px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
                >
                    Закрыть
                </button>
            </div>
        </div>
    );

    return createPortal(modalContent, document.body);
}

export default TokenDialog;