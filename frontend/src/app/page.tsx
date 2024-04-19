import { Input } from "@/components/ui/input";
import { ranHash } from "@/lib/small";
import { redirect } from "next/navigation";

export default function Home() {
  async function onsubmit(data: FormData) {
    "use server";
    let room_id = data.get("room_id")?.toString().replaceAll(" ", "-");

    return redirect(`/room/${room_id}`);
  }

  return (
    <div className="px-4 py-3 border-black border-2 m-2">
      <h1>Noob Chat</h1>
      <p>Enter in a room id, it can be anything :)</p>
      <a href="https://youtube.com">YouTube</a>
      <div>
        <form action={onsubmit} className="pt-2">
          <Input name="room_id" defaultValue={ranHash(8)} />
        </form>
      </div>
    </div>
  );
}
