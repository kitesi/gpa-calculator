import AddClassForm from "@/app/ui/AddClassForm";
export default function NewClass() {
    return <AddClassForm credits={4} year={new Date().getFullYear()} />;
}
