<?php
/**
  * @desc Testfall 1 und Testfall 2
  * @needs $_POST['object']
*/
include ('object.php'); // Objekt-Klasse für Datenstrukturierung

switch($_SERVER['REQUEST_METHOD']) {
    case "GET": // Testfall 1
        $tmp = new Object();
        $tmp->setId(substr($_SERVER['PATH_INFO'], 1));
        $tmp->setName(substr(str_shuffle("Loremipsumdolorsitametconsetetursadipscingelitr"), 0, rand(1, 20)));
        header("Content-Type: application/json");
        header("HTTP/1.1 200 OK");
    #    header("Location: /object/".$tmp->UniqId);
        echo $tmp->getJson();
    break;
    case "POST": // Testfall 2
        $input = json_decode($_POST['object']);
        if ($input != NULL) {
          $tmp = new Object();
          $tmp->setId(substr($_SERVER['PATH_INFO'], 1));
          $tmp->setName($input->Name);
          header("Content-Type: application/json");
          header("HTTP/1.1 201 Created");
          header("Location: /object/".$tmp->UniqId);
          echo $tmp->getJson();
        } else {
          header("Content-Type: application/json");
          header("HTTP/1.1 400 Bad request");
          echo json_encode(NULL);
        }
    break;
    default:
        header("HTTP/1.1 400 Bad Request");
        echo json_encode(NULL);
    break;
}
exit;
?>
