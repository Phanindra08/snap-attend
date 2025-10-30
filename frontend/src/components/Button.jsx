export default function Button({ children, type="button", onClick, className="" }) {
  return (
    <button
      type={type}
      onClick={onClick}
      className={`w-full rounded px-3 py-2 font-semibold border hover:opacity-90 ${className}`}
    >
      {children}
    </button>
  );
}
