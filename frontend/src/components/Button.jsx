export default function Button({ children, disabled, ...rest }) {
  return (
    <button className="btn" disabled={disabled} {...rest}>
      {children}
    </button>
  );
}
