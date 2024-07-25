export default function Error(props: {
    message?: string;
    children?: React.ReactNode;
}) {
    return (
        <div className="p-8">
            <h1 className="md:leading-14 mb-4 border-b-2 pb-4 text-3xl font-extrabold leading-9 tracking-tight text-gray-100 sm:text-4xl sm:leading-10 md:text-5xl">
                Error!
            </h1>
            <p className="text-white md:text-xl">
                {props.message}
                {props.children}
            </p>
        </div>
    );
}
