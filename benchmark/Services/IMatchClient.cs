public interface IMatchClient
{
    Task<JoinResponse> JoinQueue(JoinRequest request);
    Task<StatusResponse> GetStatus(string tikeckId);
}
