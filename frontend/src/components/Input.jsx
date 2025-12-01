export default function Input({
  label,
  type = "text",
  value,
  onChange,
  ...rest
}) {
  return (
    <label className="input-label">
      <span>{label}</span>
      <input
        className="input-control"
        type={type}
        value={value}
        onChange={onChange}
        {...rest}
      />
    </label>
  );
}
