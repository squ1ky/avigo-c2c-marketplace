import { useEffect, useMemo, useRef, useState } from 'react';
import { suggestListingTags } from '../../services/api';
import '../../styles/components.css';

const MIN_TITLE_LENGTH = 3;

function TagSuggestions({ title, categoryId, tags = [], onChange }) {
    const [suggestions, setSuggestions] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const requestIdRef = useRef(0);

    useEffect(() => {
        const normalizedTitle = title.trim();
        if (normalizedTitle.length < MIN_TITLE_LENGTH) {
            setSuggestions([]);
            setError('');
            setLoading(false);
            return;
        }

        const requestId = requestIdRef.current + 1;
        requestIdRef.current = requestId;

        const timeoutId = window.setTimeout(async () => {
            setLoading(true);
            setError('');

            try {
                const result = await suggestListingTags({
                    title: normalizedTitle,
                    category_id: categoryId || undefined,
                });

                if (requestIdRef.current === requestId) {
                    setSuggestions(result);
                }
            } catch (err) {
                if (requestIdRef.current === requestId) {
                    console.error(err);
                    setSuggestions([]);
                    setError(err.message || 'Не удалось подобрать теги');
                }
            } finally {
                if (requestIdRef.current === requestId) {
                    setLoading(false);
                }
            }
        }, 650);

        return () => window.clearTimeout(timeoutId);
    }, [title, categoryId]);

    const selected = useMemo(() => new Set(tags), [tags]);
    const availableSuggestions = suggestions.filter(tag => !selected.has(tag));

    const addTag = (tag) => {
        if (selected.has(tag)) return;
        onChange([...tags, tag]);
    };

    const removeTag = (tag) => {
        onChange(tags.filter(item => item !== tag));
    };

    return (
        <div className="tag-suggestions">
            <div className="tag-row tag-row-selected">
                {tags.map(tag => (
                    <button
                        key={tag}
                        type="button"
                        className="tag-chip tag-chip-selected"
                        onClick={() => removeTag(tag)}
                        title="Удалить тег"
                    >
                        {tag}
                        <span aria-hidden="true">×</span>
                    </button>
                ))}
            </div>

            <div className="tag-row">
                {availableSuggestions.map(tag => (
                    <button
                        key={tag}
                        type="button"
                        className="tag-chip"
                        onClick={() => addTag(tag)}
                    >
                        {tag}
                    </button>
                ))}
                {loading && <span className="tag-status">Подбираю...</span>}
            </div>

            {error && <div className="form-error">{error}</div>}
        </div>
    );
}

export default TagSuggestions;
