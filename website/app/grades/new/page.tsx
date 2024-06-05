import AddEditClassForm from "@/app/ui/AddEditClassForm";
export default function NewClass() {
    return (
        <AddEditClassForm
            credits={4}
            year={new Date().getFullYear()}
            className=""
            desiredGrade=""
            gradeSections={[]}
            recievedGrade=""
            semester=""
            editing={false}
        />
    );
}
