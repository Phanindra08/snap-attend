export default function Input({ label, type="text", value, onChange, placeholder, name, required }) {
  return (
    <label className="block mb-3">
      <span className="block mb-1 text-sm font-medium">{label}</span>
      <input
        className="w-full border rounded px-3 py-2 outline-none focus:ring"
        type={type}
        name={name}
        value={value}
        onChange={onChange}
        placeholder={placeholder}
        required={required}
      />
    </label>
  );
}
