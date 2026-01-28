public record GetPlayerResponse(int Id, string UserName, int Level, int Experience, int Currency)
{
    public PlayerInfo ToPlayerInfoSurvivor()
    {
        return new(Id.ToString(), UserName, "survivor");
    }

    public PlayerInfo ToPlayerInfoKiller()
    {
        return new(Id.ToString(), UserName, "killer");
    }
}


public record PlayerInfo(string Id, string Name, string Role);


