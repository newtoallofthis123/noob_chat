export default function RoomId({ params }: { params: { id: string } }) {
  return (
    <div>
      <p>{params.id}</p>
    </div>
  );
}
