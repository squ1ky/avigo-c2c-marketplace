import { useState, useRef, useEffect } from 'react';
import { uploadMedia } from '../../services/api';
import toast from 'react-hot-toast';

function ImageUploader({ mediaIds, setMediaIds, initialPreviews = [] }) {
    const [previews, setPreviews] = useState(initialPreviews);
    const [uploading, setUploading] = useState(false);
    const fileInputRef = useRef(null);

    useEffect(() => {
        if (initialPreviews.length > 0 && previews.length === 0) {
            setPreviews(initialPreviews);
        }
    }, [initialPreviews]);

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

                const response = await uploadMedia(file);

                newMediaIds.push(response.id);
                newPreviews.push({
                    id: response.id,
                    url: response.url || response.file_url
                });
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
        <div className="image-uploader" style={{width: '100%'}}>
            <div
                className="previews-grid"
                style={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fill, minmax(110px, 140px))',
                    gap: '12px',
                    marginBottom: '10px'
                }}
            >
                {previews.map((item, index) => (
                    <div
                        key={item.id}
                        className="preview-item"
                        style={{
                            position: 'relative',
                            width: '100%',
                            paddingTop: '100%',
                            borderRadius: '8px',
                            overflow: 'hidden',
                            border: '1px solid #e0e0e0',
                            backgroundColor: '#f9f9f9'
                        }}
                    >
                        <img
                            src={item.url}
                            alt={`Preview ${index}`}
                            style={{
                                position: 'absolute',
                                top: 0, left: 0,
                                width: '100%', height: '100%',
                                objectFit: 'cover'
                            }}
                        />
                        <button
                            type="button"
                            className="remove-btn"
                            onClick={() => handleRemove(item.id)}
                            title="Удалить фото"
                            style={{
                                position: 'absolute',
                                top: '4px',
                                right: '4px',
                                width: '22px',
                                height: '22px',
                                borderRadius: '50%',
                                border: 'none',
                                background: 'rgba(255, 255, 255, 0.95)',
                                color: '#ff4d4f',
                                cursor: 'pointer',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                fontSize: '16px',
                                lineHeight: '1',
                                boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
                                transition: 'transform 0.2s',
                                zIndex: 2
                            }}
                            onMouseOver={(e) => e.currentTarget.style.transform = 'scale(1.1)'}
                            onMouseOut={(e) => e.currentTarget.style.transform = 'scale(1)'}
                        >
                            ×
                        </button>
                    </div>
                ))}

                <div
                    className={`upload-btn ${uploading ? 'disabled' : ''}`}
                    onClick={() => !uploading && fileInputRef.current?.click()}
                    style={{
                        position: 'relative',
                        width: '100%',
                        paddingTop: '100%',
                        borderRadius: '8px',
                        border: '2px dashed #3f8cff',
                        backgroundColor: '#f0f7ff',
                        cursor: uploading ? 'not-allowed' : 'pointer',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        transition: 'background 0.2s'
                    }}
                    onMouseOver={(e) => !uploading && (e.currentTarget.style.backgroundColor = '#e1effe')}
                    onMouseOut={(e) => !uploading && (e.currentTarget.style.backgroundColor = '#f0f7ff')}
                >
                    <div style={{
                        position: 'absolute',
                        top: 0, left: 0, right: 0, bottom: 0,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: '#3f8cff'
                    }}>
                        {uploading ? (
                            <span style={{fontSize: '0.8rem', color: '#666'}}>Загрузка...</span>
                        ) : (
                            <>
                                <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                    <line x1="12" y1="5" x2="12" y2="19"></line>
                                    <line x1="5" y1="12" x2="19" y2="12"></line>
                                </svg>
                                <span style={{fontSize: '0.75rem', marginTop: '4px'}}>Добавить</span>
                            </>
                        )}
                    </div>
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

            <p className="hint" style={{fontSize: '0.8rem', color: '#999', margin: '0 0 5px 0'}}>
                До 10 МБ (JPG, PNG)
            </p>
        </div>
    );
}

export default ImageUploader;
