<?php
/**
  * @desc Testfall 4
  * @needs $_POST['object']
*/
include ('object.php'); // Objekt-Klasse für Datenstrukturierung

switch($_SERVER['REQUEST_METHOD']) {
    case "POST":
        $input = json_decode($_POST['object']);
        if ($input != NULL) {
          $tmp = new Object();
          $tmp->setId(substr($_SERVER['PATH_INFO'], 1));
          $tmp->setName($input->Name);

          $dbHandle = mysql_connect ("localhost", "root", "pass") or die ("keine dbHandle möglich.");
          mysql_select_db("express") or die ("Die Datenbank existiert nicht.");

          $InsertSql = "INSERT INTO tbl_person (Id, Name) VALUES('".$tmp->Id."', '".$tmp->Name."')";
          $result = mysql_query($InsertSql);

          $SelectSql = "SELECT UniqId FROM tbl_person WHERE UniqId = ".mysql_insert_id()."";
          $SelectResult = mysql_query($SelectSql);
          $ResultRow = mysql_fetch_object($SelectResult);
          mysql_close($dbHandle);

          $tmp->UniqId = $ResultRow->UniqId;

          if ($result) {
            header("Content-Type: application/json");
            header("HTTP/1.1 201 Created");
            header("Location: /object/".$tmp->UniqId);
            echo $tmp->getJson();
          } else {
            header("Content-Type: application/json");
            header("HTTP/1.1 400 Bad request");
            echo json_encode(NULL);
          }
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
