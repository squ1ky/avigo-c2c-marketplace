import { useState, useRef } from 'react';
import { uploadMedia } from '../../services/api';
import toast from 'react-hot-toast';

function ImageUploader({ mediaIds, setMediaIds, initialPreviews = [] }) {
    const [previews, setPreviews] = useState(initialPreviews);
    const [uploading, setUploading] = useState(false);
    const fileInputRef = useRef(null);

    const handleFileSelect = async (e) => {
        const files = Array.from(e.target.files);
        if (files.length === 0) return;

        setUploading(true);
        const newMediaIds = [...mediaIds];
        const newPreviews = [...previews];

        try {
            for (const file of files) {
                if (file.size > 10 * 1024 * 1024) {
                    toast.error(`Файл ${file.name} слишком большой (>10MB)`);
                    continue;
                }

                const response = await uploadMedia(file); // { id, url }

                newMediaIds.push(response.id);
                newPreviews.push({ id: response.id, url: response.url });
            }

            setMediaIds(newMediaIds);
            setPreviews(newPreviews);
            toast.success('Фото загружены');
        } catch (error) {
            console.error(error);
            toast.error('Ошибка загрузки фото');
        } finally {
            setUploading(false);
            if (fileInputRef.current) fileInputRef.current.value = '';
        }
    };

    const handleRemove = (idToRemove) => {
        const newMediaIds = mediaIds.filter(id => id !== idToRemove);
        const newPreviews = previews.filter(p => p.id !== idToRemove);

        setMediaIds(newMediaIds);
        setPreviews(newPreviews);
    };

    return (
        <div className="image-uploader">
            <div className="previews-grid">
                {previews.map((item, index) => (
                    <div key={item.id} className="preview-item">
                        <img src={item.url} alt={`Preview ${index}`} />
                        <button
                            type="button"
                            className="remove-btn"
                            onClick={() => handleRemove(item.id)}
                        >
                            ×
                        </button>
                    </div>
                ))}

                <div
                    className={`upload-btn ${uploading ? 'disabled' : ''}`}
                    onClick={() => !uploading && fileInputRef.current?.click()}
                >
                    {uploading ? (
                        <span className="loader">...</span>
                    ) : (
                        <>
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                <path d="M12 5v14M5 12h14" />
                            </svg>
                            <span>Добавить фото</span>
                        </>
                    )}
                </div>
            </div>

            <input
                type="file"
                ref={fileInputRef}
                onChange={handleFileSelect}
                multiple
                accept="image/*"
                style={{ display: 'none' }}
            />

            <p className="hint">Перетащите фото или нажмите кнопку. До 10 МБ.</p>
        </div>
    );
}

export default ImageUploader;
