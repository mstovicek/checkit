<?php

if (!isset($_GET['uuid'])) {
    echo '{"return_code": "0", "error": "invalid uuid"}';
    return;
}

$output = exec(
    sprintf(
        'php-cs-fixer fix -vvv --diff --dry-run --format=json --using-cache=no /files/%s 2>&1',
        $_GET['uuid']
    )
);

echo sprintf(
    '{"return_code": "1", "output": %s}',
    $output ?: '""'
);
