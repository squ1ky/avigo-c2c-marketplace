import { useState, useEffect } from 'react';
import { getCategories } from '../../services/api';
import toast from 'react-hot-toast';

const flattenCategories = (categories, level = 0, result = []) => {
    for (const cat of categories) {
        result.push({
            id: cat.id,
            name: cat.name,
            level: level,
            disabled: cat.children && cat.children.length > 0
        });

        if (cat.children && cat.children.length > 0) {
            flattenCategories(cat.children, level + 1, result);
        }
    }
    return result;
};

function CategorySelect({ value, onChange }) {
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const load = async () => {
            try {
                const tree = await getCategories();
                const flatList = flattenCategories(tree);
                setCategories(flatList);
            } catch (err) {
                console.error(err);
                toast.error('Не удалось загрузить категории');
            } finally {
                setLoading(false);
            }
        };
        load();
    }, []);

    if (loading) return <div className="form-input" style={{color: '#999'}}>Загрузка категорий...</div>;

    return (
        <select
            className="form-input"
            value={value}
            onChange={(e) => onChange(e.target.value)}
            required
            style={{
                fontWeight: value ? '500' : 'normal',
                cursor: 'pointer'
            }}
        >
            <option value="">Выберите категорию</option>
            {categories.map((cat) => (
                <option
                    key={cat.id}
                    value={cat.id}
                    disabled={cat.disabled}
                    style={{
                        fontWeight: cat.level === 0 ? 'bold' : 'normal',
                        color: cat.disabled ? '#999' : '#000'
                    }}
                >
                    {'\u00A0'.repeat(cat.level * 4)}
                    {cat.name}
                </option>
            ))}
        </select>
    );
}

export default CategorySelect;
