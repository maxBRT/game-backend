namespace player_service.Exceptions;

public class PlayerNotFoundException(string message) : Exception(message)
{
}
